// internal/services/user_service.go
package services

import (
	"SDGEStreaming/internal/models"
	"SDGEStreaming/internal/repositories"
	"SDGEStreaming/internal/security"
	"SDGEStreaming/internal/utils"
	"errors"
	"fmt"
	"time"
)

type UserService struct {
	userRepo         repositories.UserRepo
	subscriptionRepo repositories.SubscriptionRepo
}

func NewUserService(userRepo repositories.UserRepo, subscriptionRepo repositories.SubscriptionRepo) *UserService {
	return &UserService{userRepo: userRepo, subscriptionRepo: subscriptionRepo}
}

func (s *UserService) Register(name string, age int, email, password string, isAdmin bool) (*models.User, error) {
	if !utils.IsValidName(name) {
		return nil, fmt.Errorf("nombre inválido")
	}
	if age < 13 || age > 120 {
		return nil, fmt.Errorf("edad debe estar entre 13 y 120 años")
	}
	if !utils.IsValidEmail(email) {
		return nil, fmt.Errorf("email inválido")
	}
	if !utils.IsValidPassword(password) {
		return nil, fmt.Errorf("contraseña debe tener al menos 6 caracteres")
	}

	if existing, _ := s.userRepo.FindByEmail(email); existing != nil {
		return nil, fmt.Errorf("el email ya está registrado")
	}

	hashedPass, err := security.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("error al procesar la contraseña")
	}

	// Clasificación de edad automática
	ageRating := classifyAge(age)

	now := time.Now()
	user := &models.User{
		Name:         name,
		Age:          age,
		Email:        email,
		PasswordHash: hashedPass,
		PlanID:       1, // plan Free por defecto
		AgeRating:    ageRating,
		IsAdmin:      isAdmin,
		CreatedAt:    now,
		LastLogin:    now,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("no se pudo crear la cuenta de usuario: %w", err)
	}

	return user, nil
}

// classifyAge clasifica la edad del usuario en categorías
func classifyAge(age int) string {
	switch {
	case age < 13:
		return "Niño"
	case age < 18:
		return "Adolescente"
	default:
		return "Adulto"
	}
}

// Login existing signature kept
func (s *UserService) Login(email, password string) (*models.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil || user == nil {
		return nil, errors.New("email o contraseña incorrectos")
	}

	if !security.CheckPasswordHash(password, user.PasswordHash) {
		return nil, errors.New("email o contraseña incorrectos")
	}

	return user, nil
}

// GetByID retrieves a user by ID.
func (s *UserService) GetByID(id int) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

// GetAllUsers para el admin
func (s *UserService) GetAllUsers() ([]models.User, error) {
	users, err := s.userRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("no se pudo obtener los usuarios: %w", err)
	}
	return users, nil
}

// UpdateUserPlan actualiza el plan (usado por main)
func (s *UserService) UpdateUserPlan(userID, planID int) error {
	// actualizar tabla users
	if err := s.userRepo.UpdatePlan(userID, planID); err != nil {
		return fmt.Errorf("no se pudo actualizar plan en users: %w", err)
	}
	return nil
}

// GetDefaultPaymentMethod
func (s *UserService) GetDefaultPaymentMethod(userID int) (*models.PaymentMethod, error) {
	return s.userRepo.GetDefaultPaymentMethod(userID)
}
