package service

import (
	"belajar-backend-golang/internal/config"
	"belajar-backend-golang/internal/model"
	"belajar-backend-golang/internal/repository"
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Repo *repository.UserRepository
	Cfg  *config.Config
}

func NewUserService(r *repository.UserRepository, c *config.Config) *UserService {
	return &UserService{
		Repo: r,
		Cfg:  c,
	}
}

func (s *UserService) Register(ctx context.Context, username, email, password string) (model.User, error) {
	existing, err := s.Repo.GetUserByEmail(ctx, email)
	if err != nil && err.Error() != "no rows in result set" {
		return model.User{}, err
	}
	if existing.ID != 0 {
		return model.User{}, errors.New("Email is already used")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, err
	}
	u := model.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}
	neweUser, err := s.Repo.CreateUser(ctx, u)
	if err != nil {
		return model.User{}, err
	}
	return neweUser, nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (model.User, string, error) {
	findUser, err := s.Repo.GetUserByEmail(ctx, email)
	if err != nil {
		return model.User{}, "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(findUser.Password), []byte(password))
	if err != nil {
		return model.User{}, "", errors.New("Invalid password")
	}

	token, err := s.GenerateJWTToken(findUser.ID)
	if err != nil {
		return model.User{}, "", err
	}

	return findUser, token, nil
}

func (s *UserService) GenerateJWTToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.Cfg.JWTSECRET))
}
