package application

import (
	"errors"

	"github.com/user/mangahub/internal/domain"
	"github.com/user/mangahub/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserExists   = errors.New("username already exists")
	ErrInvalidLogin = errors.New("invalid username or password")
)

type AuthService interface {
	Register(username, password, role string) (*models.User, error)
	Login(username, password string) (*models.User, error)
}

type authService struct {
	userRepo domain.UserRepository
}

// NewAuthService strictly instantiates isolated Authentication logic flows.
func NewAuthService(repo domain.UserRepository) AuthService {
	return &authService{userRepo: repo}
}

func (s *authService) Register(username, password, role string) (*models.User, error) {
	if existing, _ := s.userRepo.FindByUsername(username); existing != nil {
		return nil, ErrUserExists
	}

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &models.User{
		Username:     username,
		PasswordHash: string(hashBytes),
		Role:         role,
	}

	if err := s.userRepo.Save(u); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *authService) Login(username, password string) (*models.User, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil || user == nil {
		return nil, ErrInvalidLogin
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidLogin
	}

	return user, nil
}
