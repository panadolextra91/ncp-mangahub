package database

import (
	"database/sql"
	"github.com/user/mangahub/pkg/models"
)

type sqliteUserRepo struct {
	db *sql.DB
}

func NewSqliteUserRepository(db *sql.DB) *sqliteUserRepo {
	return &sqliteUserRepo{db: db}
}

func (r *sqliteUserRepo) Save(user *models.User) error {
	query := `INSERT INTO users (username, password_hash, role) VALUES (?, ?, ?)`
	res, err := r.db.Exec(query, user.Username, user.PasswordHash, user.Role)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = int(id)
	return nil
}

func (r *sqliteUserRepo) FindByUsername(username string) (*models.User, error) {
	query := `SELECT id, username, password_hash, role FROM users WHERE username = ?`
	row := r.db.QueryRow(query, username)
	
	user := &models.User{}
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}
