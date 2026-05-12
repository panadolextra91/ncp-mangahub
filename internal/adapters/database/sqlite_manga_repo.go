package database

import (
	"database/sql"
	"strings"

	"github.com/user/mangahub/internal/domain"
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

// SearchWithFilters runs an advanced multi-criteria query against the mangas
// table. All filter fields are optional and combinable as AND across different
// types (query AND genres AND status); genres themselves are OR'd together.
// Sort defaults to id DESC ("recent"); "title" sorts LOWER(title) ASC; any
// other sortBy value falls back to recent. No schema migration — strictly
// works against existing columns.
func (r *sqliteMangaRepo) SearchWithFilters(f domain.SearchFilters) ([]*models.Manga, error) {
	var (
		clauses []string
		args    []interface{}
	)

	if f.Query != "" {
		clauses = append(clauses, "(LOWER(title) LIKE LOWER(?) OR LOWER(author) LIKE LOWER(?))")
		like := "%" + f.Query + "%"
		args = append(args, like, like)
	}

	// Defensive cap at 10 — the HTTP handler also caps, but belt-and-suspenders.
	genres := f.Genres
	if len(genres) > 10 {
		genres = genres[:10]
	}
	if len(genres) > 0 {
		var genreClauses []string
		for _, g := range genres {
			g = strings.TrimSpace(g)
			if g == "" {
				continue
			}
			// Quoted-token form prevents substring collisions ("Action" vs "Reaction").
			genreClauses = append(genreClauses, "genres LIKE ?")
			args = append(args, "%\""+g+"\"%")
		}
		if len(genreClauses) > 0 {
			clauses = append(clauses, "("+strings.Join(genreClauses, " OR ")+")")
		}
	}

	if f.Status != "" {
		clauses = append(clauses, "status = ?")
		args = append(args, f.Status)
	}

	where := "1=1"
	if len(clauses) > 0 {
		where = strings.Join(clauses, " AND ")
	}

	orderBy := "id DESC"
	if f.SortBy == "title" {
		orderBy = "LOWER(title) ASC"
	}
	// Any other value (including "recent", "", or unknown) → id DESC.

	query := "SELECT id, title, author, genres, status, total_chapters, description, created_at FROM mangas WHERE " + where + " ORDER BY " + orderBy
	rows, err := r.db.Query(query, args...)
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
