package service

import (
	"car-showroom-api/internal/entity"
	"car-showroom-api/internal/repository"
	"car-showroom-api/pkg/utils"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Phone    string `json:"phone" binding:"required"`
	Role     string `json:"role"` // Opsional saat register, default: buyer
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthService interface {
	Register(input RegisterInput) (*entity.User, error)
	Login(input LoginInput) (string, *entity.User, error)
	GetProfile(userID uint) (*entity.User, error)
}

type authService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &authService{repo}
}

func (s *authService) Register(input RegisterInput) (*entity.User, error) {
	existingUser, _ := s.repo.FindByEmail(input.Email)
	if existingUser.ID != 0 {
		return nil, errors.New("email sudah terdaftar")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("gagal memproses password")
	}

	role := "buyer"
	if input.Role == "admin" || input.Role == "sales" {
		role = input.Role
	}

	user := &entity.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
		Phone:    input.Phone,
		Role:     role,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(input LoginInput) (string, *entity.User, error) {
	user, err := s.repo.FindByEmail(input.Email)
	if err != nil || user.ID == 0 {
		return "", nil, errors.New("email atau password salah")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return "", nil, errors.New("email atau password salah")
	}

	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", nil, errors.New("gagal membuat token autentikasi")
	}

	return token, user, nil
}

func (s *authService) GetProfile(userID uint) (*entity.User, error) {
	return s.repo.FindByID(userID)
}
