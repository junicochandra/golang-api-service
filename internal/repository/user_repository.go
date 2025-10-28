package repository

import (
	"database/sql"

	"github.com/junicochandra/golang-api-service/internal/config/database"
	"github.com/junicochandra/golang-api-service/internal/entity"
)

type UserRepository interface {
	FindAll() ([]entity.User, error)
	FindById(id uint64) (*entity.User, error)
	FindByEmail(email string) (bool, error)
	Create(user *entity.User) (uint64, error)
	Update(user *entity.User) (int64, error)
	Delete(id uint64) (int64, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: database.DB}
}

func (repo *userRepository) FindAll() ([]entity.User, error) {
	rows, err := repo.db.Query("SELECT id, name, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var user entity.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (repo *userRepository) FindById(id uint64) (*entity.User, error) {
	var user entity.User
	row := repo.db.QueryRow("SELECT id, name, email FROM users WHERE id = ?", id)

	if err := row.Scan(&user.ID, &user.Name, &user.Email); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (repo *userRepository) FindByEmail(email string) (bool, error) {
	var exists bool
	err := repo.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", email).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (repo *userRepository) Create(user *entity.User) (uint64, error) {
	res, err := repo.db.Exec("INSERT INTO users SET name = ?, email = ?, password = ?", user.Name, user.Email, user.Password)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(id), nil
}

func (repo *userRepository) Update(user *entity.User) (int64, error) {
	res, err := repo.db.Exec("UPDATE users SET name = ?, email = ? WHERE id = ?", user.Name, user.Email, user.ID)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func (repo *userRepository) Delete(id uint64) (int64, error) {
	res, err := repo.db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}
