package service

import (
	"dailyworkerroster/middleware"
	"dailyworkerroster/model"
	"dailyworkerroster/repository"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserServiceItf interface {
	SignUp(user *model.User) (int64, error)
	Login(identifier, password string) (*model.User, error)
	GetAllWorkers() ([]*model.User, error)
	GetWorkerByID(workerID int64) (*model.User, error)
}

type UserService struct {
	UserRepo repository.UserRepoItf
}

func NewUserService(
	userRepo repository.UserRepoItf) UserServiceItf {
	return &UserService{
		UserRepo: userRepo,
	}
}

func (s *UserService) SignUp(user *model.User) (int64, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	user.Password = string(hashedPassword)
	return s.UserRepo.SignUp(user)
}

func (s *UserService) Login(identifier, password string) (*model.User, error) {
	user, err := s.UserRepo.Login(identifier)
	if err != nil {
		return nil, errors.New("failed to get login credentials")
	}
	fmt.Println("halo", user.Password, " ", password)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid username/email or password")
	}

	user.Password = ""
	user.JWTToken, err = middleware.GenerateJWT(user.ID, user.Name, user.Role)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetAllWorkers() ([]*model.User, error) {
	return s.UserRepo.GetUsersByRole(model.ROLE_WORKER)
}

func (s *UserService) GetWorkerByID(workerID int64) (*model.User, error) {
	user, err := s.UserRepo.GetUserByID(workerID)
	if err != nil {
		return nil, err
	}
	if user.Role != model.ROLE_WORKER {
		return nil, errors.New("user is not a worker")
	}
	user.Password = ""
	return user, nil
}
