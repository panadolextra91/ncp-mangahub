package application_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/user/mangahub/internal/application"
	"github.com/user/mangahub/internal/domain"
	"github.com/user/mangahub/internal/eventbus"
	"github.com/user/mangahub/pkg/models"
)

// --- MOCKS ---

type mockUserRepo struct {
	users map[string]*models.User
}

func (m *mockUserRepo) Save(u *models.User) error {
	m.users[u.Username] = u
	return nil
}

func (m *mockUserRepo) FindByUsername(username string) (*models.User, error) {
	u, ok := m.users[username]
	if !ok {
		return nil, errors.New("not found")
	}
	return u, nil
}

type mockMangaRepo struct {
	savedCount int
}

func (m *mockMangaRepo) Save(manga *models.Manga) error {
	m.savedCount++
	return nil
}
func (m *mockMangaRepo) FindByID(id int) (*models.Manga, error) { return nil, nil }
func (m *mockMangaRepo) List() ([]*models.Manga, error)          { return nil, nil }
func (m *mockMangaRepo) Search(q string) ([]*models.Manga, error) { return nil, nil }
func (m *mockMangaRepo) SearchWithFilters(f domain.SearchFilters) ([]*models.Manga, error) {
	return nil, nil
}

type mockProgressRepo struct{}

func (m *mockProgressRepo) Save(p *models.UserProgress) error                                     { return nil }
func (m *mockProgressRepo) FindByUserAndManga(userID, mangaID int) (*models.UserProgress, error) { return nil, nil }
func (m *mockProgressRepo) GetByUserID(userID int) ([]*models.UserProgress, error)               { return nil, nil }

// --- TESTS ---

func TestAuthService_Logic(t *testing.T) {
	repo := &mockUserRepo{users: make(map[string]*models.User)}
	authSvc := application.NewAuthService(repo)

	// 1. Success Register
	user, err := authSvc.Register("thuy", "supersecret", "admin")
	assert.NoError(t, err)
	assert.Equal(t, "thuy", user.Username)
	assert.NotEqual(t, "supersecret", user.PasswordHash) // Strict Bcrypt check

	// 2. Duplicate Check
	_, err = authSvc.Register("thuy", "another", "user")
	assert.Equal(t, application.ErrUserExists, err)

	// 3. Login Success
	loginUser, err := authSvc.Login("thuy", "supersecret")
	assert.NoError(t, err)
	assert.Equal(t, "admin", loginUser.Role)

	// 4. Bad Password Rejection
	_, err = authSvc.Login("thuy", "wrong")
	assert.Equal(t, application.ErrInvalidLogin, err)
}

func TestMangaService_CreateRoleBlocks(t *testing.T) {
	repo := &mockMangaRepo{}
	bus := eventbus.NewEventBus(10)
	ch := bus.Subscribe("manga.new")

	mangaSvc := application.NewMangaService(repo, bus)
	newManga := &models.Manga{Title: "Naruto"}

	// 1. Standard user role strictly rejected (Simulation)
	err := mangaSvc.CreateManga("user", newManga)
	assert.Equal(t, application.ErrUnauthorizedCreate, err)
	assert.Equal(t, 0, repo.savedCount)

	// 2. Admin role accepted and native event published accurately
	err = mangaSvc.CreateManga("admin", newManga)
	assert.NoError(t, err)
	assert.Equal(t, 1, repo.savedCount)

	select {
	case ev := <-ch:
		assert.Equal(t, "manga.new", ev.Topic)
	case <-time.After(1 * time.Second):
		t.Fatal("EventBus omitted Manga creation payload!")
	}
}

func TestProgressService_Logic(t *testing.T) {
	repo := &mockProgressRepo{}
	bus := eventbus.NewEventBus(10)
	ch := bus.Subscribe("progress.updated")

	svc := application.NewProgressService(repo, bus)
	err := svc.UpdateProgress(&models.UserProgress{UserID: 1, MangaID: 2})
	assert.NoError(t, err)

	select {
	case ev := <-ch:
		assert.Equal(t, "progress.updated", ev.Topic)
	case <-time.After(1 * time.Second):
		t.Fatal("EventBus omitted progress update payload!")
	}
}
