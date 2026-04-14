package domain

type Manga struct {
	ID    int
	Title string
}

type MangaRepository interface {
	Save(manga Manga) error
}
