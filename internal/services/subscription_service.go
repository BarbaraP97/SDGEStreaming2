// internal/services/subscription_service.go
package services

import (
	"SDGEStreaming/internal/models"
	"SDGEStreaming/internal/repositories"
	"fmt"
)

type SubscriptionService struct {
	subRepo  repositories.SubscriptionRepo
	userRepo repositories.UserRepo
}

func NewSubscriptionService(subRepo repositories.SubscriptionRepo, userRepo repositories.UserRepo) *SubscriptionService {
	return &SubscriptionService{subRepo: subRepo, userRepo: userRepo}
}

func (s *SubscriptionService) ProcessPayment(userID int, planID int, cardHolder, cardNumber string, expiryMonth, expiryYear, cvv int) error {
	// Validación básica de la tarjeta
	if len(cardNumber) < 13 || len(cardNumber) > 19 {
		return fmt.Errorf("número de tarjeta inválido")
	}
	if expiryMonth < 1 || expiryMonth > 12 {
		return fmt.Errorf("mes de vencimiento inválido")
	}
	if cvv < 100 || cvv > 999 {
		return fmt.Errorf("CVV inválido")
	}

	// Obtener el plan
	plan, err := s.subRepo.GetPlanByID(planID)
	if err != nil {
		return err
	}

	// Simular el procesamiento del pago
	fmt.Printf("Procesando pago de $%.2f para el plan '%s'...\n", plan.Price, plan.Name)
	fmt.Println("¡Pago aprobado!")

	// Guardar el método de pago
	last4 := cardNumber[len(cardNumber)-4:]
expirationDate := fmt.Sprintf("%02d/%04d", expiryMonth, expiryYear)

method := &models.PaymentMethod{
    UserID:        userID,
    CardHolder:    cardHolder,
    CardNumber:    cardNumber,
    ExpirationDate: expirationDate,
    CVV:           fmt.Sprintf("%03d", cvv),
    Last4:         last4,
    ExpiryMonth:   expiryMonth,
    ExpiryYear:    expiryYear,
    IsDefault:     true,
}
	if err := s.userRepo.AddPaymentMethod(method); err != nil {
		return fmt.Errorf("error al guardar el método de pago")
	}

	// Actualizar el plan del usuario
	if err := s.userRepo.UpdatePlan(userID, planID); err != nil {
		return fmt.Errorf("error al actualizar el plan")
	}

	return nil
}

func (s *SubscriptionService) GetAvailablePlans() ([]models.Plan, error) {
	return s.subRepo.GetAllPlans()
}
