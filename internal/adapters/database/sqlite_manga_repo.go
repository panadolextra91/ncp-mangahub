package database

import (
	"database/sql"
	"github.com/user/mangahub/pkg/models"
)

type sqliteMangaRepo struct {
	db *sql.DB
}

func NewSqliteMangaRepository(db *sql.DB) *sqliteMangaRepo {
	return &sqliteMangaRepo{db: db}
}

func (r *sqliteMangaRepo) Save(manga *models.Manga) error {
	query := `INSERT INTO mangas (title, author) VALUES (?, ?)`
	res, err := r.db.Exec(query, manga.Title, manga.Author)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	manga.ID = int(id)
	return nil
}

func (r *sqliteMangaRepo) FindByID(id int) (*models.Manga, error) {
	query := `SELECT id, title, author, created_at FROM mangas WHERE id = ?`
	row := r.db.QueryRow(query, id)
	
	m := &models.Manga{}
	err := row.Scan(&m.ID, &m.Title, &m.Author, &m.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return m, nil
}
