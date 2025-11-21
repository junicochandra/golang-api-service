package bootstrap

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/junicochandra/golang-api-service/internal/domain/entity"
	"github.com/junicochandra/golang-api-service/internal/infrastructure/config/database"
	"github.com/junicochandra/golang-api-service/internal/infrastructure/repository"
	rabbitmq "github.com/junicochandra/golang-api-service/internal/infrastructure/service/rabbitmq"
	worker "github.com/junicochandra/golang-api-service/internal/infrastructure/service/rabbitmq/worker"
	"github.com/junicochandra/golang-api-service/internal/router"
)

func Run() {
	// DB init
	database.Connect()
	db := database.DB
	if err := db.AutoMigrate(&entity.User{}); err != nil {
		log.Fatal("migrate error: ", err)
	}

	accountRepo := repository.NewAccountRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	// RabbitMQ init
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}
	rabbitSvc, err := rabbitmq.New(rabbitURL)
	if err != nil {
		log.Fatalf("rabbit connect error: %v", err)
	}
	defer rabbitSvc.Close()

	// Declare topology
	err = rabbitmq.DeclareTopology(rabbitSvc, rabbitmq.TopologyConfig{
		Exchange:   "topup.exchange",
		ExchangeTy: "direct",
		Queue:      "topup_queue",
		RoutingKey: "topup.created",
		// DLX: "topup.dlx", // activate if using DLX
	})
	if err != nil {
		log.Fatalf("declare topology error: %v", err)
	}

	// Router (adjust if the router accepts the usecase)
	r := router.SetupRouter(rabbitSvc)

	// Start worker
	logger := log.New(os.Stdout, "[topup-worker] ", log.LstdFlags)
	cons := worker.NewConsumer(rabbitSvc, transactionRepo, accountRepo, "topup_queue", logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := cons.Start(ctx); err != nil {
			logger.Fatalf("worker error: %v", err)
		}
	}()

	// Run server (non-blocking)
	serverErr := make(chan error, 1)
	go func() {
		serverErr <- r.Run(":9000")
	}()

	// Graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-sig:
		log.Printf("signal: %v, shutting down", s)
		cancel()
		// give the worker time to complete the task
		time.Sleep(2 * time.Second)
	case err := <-serverErr:
		log.Printf("server error: %v", err)
		cancel()
	}

	log.Println("bootstrap: stopped")
}
