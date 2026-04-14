package database

import (
	"database/sql"
	"github.com/user/mangahub/pkg/models"
)

type sqliteProgressRepo struct {
	db *sql.DB
}

func NewSqliteProgressRepository(db *sql.DB) *sqliteProgressRepo {
	return &sqliteProgressRepo{db: db}
}

func (r *sqliteProgressRepo) Save(progress *models.UserProgress) error {
	query := `INSERT INTO user_progress (user_id, manga_id, current_chapter, updated_at) 
			  VALUES (?, ?, ?, CURRENT_TIMESTAMP)
			  ON CONFLICT(user_id, manga_id) DO UPDATE SET
			  current_chapter = excluded.current_chapter,
			  updated_at = CURRENT_TIMESTAMP`
	_, err := r.db.Exec(query, progress.UserID, progress.MangaID, progress.CurrentChapter)
	return err
}

func (r *sqliteProgressRepo) FindByUserAndManga(userID, mangaID int) (*models.UserProgress, error) {
	query := `SELECT user_id, manga_id, current_chapter, updated_at FROM user_progress WHERE user_id = ? AND manga_id = ?`
	row := r.db.QueryRow(query, userID, mangaID)
	
	p := &models.UserProgress{}
	err := row.Scan(&p.UserID, &p.MangaID, &p.CurrentChapter, &p.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return p, nil
}
