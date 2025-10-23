package service

import (
	"context"
	"errors"

	"user-service/internal/model"
	"user-service/internal/repository"
	"user-service/internal/util"

	pb "user-service/user"

	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	Repo repository.UserRepository
}

var (
	ErrUsernameTaken      = errors.New("username already taken")
	ErrEmailRegistered    = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func (s *UserService) RegisterUser(username, email, password, firstName, lastName string) error {
	_, err := s.Repo.GetUserByUsername(username)
	if err == nil {
		return ErrUsernameTaken
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	_, err = s.Repo.GetUserByEmail(email)
	if err == nil {
		return ErrEmailRegistered
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := model.User{
		FirstName:    firstName,
		LastName:     lastName,
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}
	return s.Repo.CreateUser(&user)
}

func (s *UserService) LoginUser(username, password string) (*model.User, error) {
	user, err := s.Repo.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	return user, nil
}

func (s *UserService) AuthenticateUser(ctx context.Context, req *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {
	user, err := s.Repo.GetUserByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(grpcCodes.Unauthenticated, "invalid credentials")
		}
		return nil, status.Errorf(grpcCodes.Unauthenticated, "invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, status.Errorf(grpcCodes.Unauthenticated, "invalid credentials")
	}

	token, err := util.GenerateJWT(user.Username)
	if err != nil {
		return nil, status.Errorf(grpcCodes.Unauthenticated, "invalid credentials")
	}
	return &pb.AuthenticateResponse{Token: token}, nil
}

func (s *UserService) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	claims, err := util.ValidateJWT(req.Token)
	if err != nil {
		return &pb.ValidateTokenResponse{Username: claims.Username, Valid: true}, nil
	}
	return &pb.ValidateTokenResponse{Valid: false}, nil
}
