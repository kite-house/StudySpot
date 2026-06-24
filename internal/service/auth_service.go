package service

import (
	"errors"
	"time"

	"studyspot/internal/domain"
	"studyspot/internal/repository"
	"studyspot/pkg/jwt"
	"studyspot/pkg/password"

	"github.com/google/uuid"
)

type AuthService struct {
	userRepo       *repository.UserRepository
	jwtSecret      string
	jwtExpireHours int
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string, jwtExpireHours int) *AuthService {
	return &AuthService{
		userRepo:       userRepo,
		jwtSecret:      jwtSecret,
		jwtExpireHours: jwtExpireHours,
	}
}

func (s *AuthService) Register(email, pass string) (*domain.User, error) {
	existingUser, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := password.Hash(pass)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		ID:        uuid.New(),
		Email:     email,
		Password:  hashedPassword,
		Role:      "user",
		CreatedAt: time.Now(),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	user.Password = ""
	return user, nil
}

func (s *AuthService) Login(email, pass string) (string, *domain.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", nil, err
	}
	if user == nil {
		return "", nil, errors.New("invalid credentials")
	}

	if !password.CheckHash(pass, user.Password) {
		return "", nil, errors.New("invalid credentials")
	}

	token, err := jwt.GenerateToken(user.ID, user.Email, user.Role, s.jwtSecret, s.jwtExpireHours)
	if err != nil {
		return "", nil, err
	}

	user.Password = ""
	return token, user, nil
}
