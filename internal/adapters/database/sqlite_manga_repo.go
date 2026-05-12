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
	query := `INSERT INTO mangas (title, author, genres, status, total_chapters, description) VALUES (?, ?, ?, ?, ?, ?)`
	res, err := r.db.Exec(query, manga.Title, manga.Author, manga.Genres, manga.Status, manga.TotalChapters, manga.Description)
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
	query := `SELECT id, title, author, genres, status, total_chapters, description, created_at FROM mangas WHERE id = ?`
	row := r.db.QueryRow(query, id)

	m := &models.Manga{}
	err := row.Scan(&m.ID, &m.Title, &m.Author, &m.Genres, &m.Status, &m.TotalChapters, &m.Description, &m.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return m, nil
}

func (r *sqliteMangaRepo) List() ([]*models.Manga, error) {
	query := `SELECT id, title, author, genres, status, total_chapters, description, created_at FROM mangas ORDER BY id DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*models.Manga
	for rows.Next() {
		m := &models.Manga{}
		if err := rows.Scan(&m.ID, &m.Title, &m.Author, &m.Genres, &m.Status, &m.TotalChapters, &m.Description, &m.CreatedAt); err != nil {
			return nil, err
		}
		results = append(results, m)
	}
	return results, nil
}

func (r *sqliteMangaRepo) Search(q string) ([]*models.Manga, error) {
	query := `SELECT id, title, author, genres, status, total_chapters, description, created_at 
	          FROM mangas 
	          WHERE LOWER(title) LIKE LOWER(?) OR LOWER(author) LIKE LOWER(?) OR LOWER(genres) LIKE LOWER(?) 
	          ORDER BY id DESC`
	searchQuery := "%" + q + "%"
	rows, err := r.db.Query(query, searchQuery, searchQuery, searchQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*models.Manga
	for rows.Next() {
		m := &models.Manga{}
		if err := rows.Scan(&m.ID, &m.Title, &m.Author, &m.Genres, &m.Status, &m.TotalChapters, &m.Description, &m.CreatedAt); err != nil {
			return nil, err
		}
		results = append(results, m)
	}
	return results, nil
}
