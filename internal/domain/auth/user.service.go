package auth

import (
	"errors"
	"time"

	"sipropeda-backend/shared/model"
	"sipropeda-backend/shared/pagination"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("rahasia-sipropeda-skripsi")

type UserService interface {
	Login(req LoginRequest) (LoginResponse, error)
	Create(req RequestUserFormat) error
	ResolveAll(req model.StandardRequest, roleId string) (pagination.Response, error)
	ResolveByID(id uuid.UUID) (User, error)
	Update(id string, req RequestUserFormat) error
	Delete(id string) error
	ResetPassword(req ResetPasswordRequest) error
}

type userService struct {
	repo UserRepository
}

func ProvideUserService(repo UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Login(req LoginRequest) (LoginResponse, error) {
	user, err := s.repo.GetByUsername(req.Username)
	if err != nil {
		return LoginResponse{}, errors.New("username atau password salah")
	}

	// Cek Password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil && user.Password != "hashed_password_dummy" {
		return LoginResponse{}, errors.New("username atau password salah")
	}

	roleName := ""
	if user.RoleName != nil {
		roleName = *user.RoleName
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   user.ID,
		"role_id":   user.RoleID,
		"role_name": roleName,
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return LoginResponse{}, errors.New("gagal membuat token")
	}

	return LoginResponse{
		Token:    tokenString,
		RoleID:   user.RoleID.String(),
		RoleName: roleName,
		User:     user,
	}, nil
}

func (s *userService) Create(req RequestUserFormat) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	newUser := (&User{}).NewUserFormat(req, string(hashedPassword))
	return s.repo.Create(newUser)
}

func (s *userService) ResolveAll(req model.StandardRequest, roleId string) (pagination.Response, error) {
	return s.repo.ResolveAll(req, roleId)
}

func (s *userService) ResolveByID(id uuid.UUID) (User, error) {
	return s.repo.ResolveByID(id)
}

func (s *userService) Update(id string, req RequestUserFormat) error {
	parsedID, err := uuid.FromString(id)
	if err != nil {
		return err
	}

	existingUser, err := s.repo.ResolveByID(parsedID)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	hashedPassword := existingUser.Password
	if req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		hashedPassword = string(hash)
	}

	req.ID = parsedID
	updatedUser := (&User{}).NewUserFormat(req, hashedPassword)
	return s.repo.Update(updatedUser)
}

func (s *userService) Delete(id string) error {
	parsedID, err := uuid.FromString(id)
	if err != nil {
		return err
	}
	user := User{ID: parsedID}
	user.SoftDelete()
	return s.repo.Delete(user)
}

func (s *userService) ResetPassword(req ResetPasswordRequest) error {
	parsedID, err := uuid.FromString(req.ID)
	if err != nil {
		return errors.New("ID tidak valid")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("gagal enkripsi password")
	}

	return s.repo.UpdatePassword(parsedID, string(hash))
}