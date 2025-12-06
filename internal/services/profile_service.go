// internal/services/profile_service.go
package services

import (
	"SDGEStreaming/internal/models"
	"SDGEStreaming/internal/repositories"
	"fmt"
	"time"
)

type ProfileService struct {
	profileRepo repositories.ProfileRepo
	userRepo    repositories.UserRepo
}

func NewProfileService(profileRepo repositories.ProfileRepo, userRepo repositories.UserRepo) *ProfileService {
	return &ProfileService{
		profileRepo: profileRepo,
		userRepo:    userRepo,
	}
}

func (s *ProfileService) CreateProfile(userID int, name, profileType string, isMain bool) (*models.Profile, error) {
	ageRating := "R"
	switch profileType {
	case "kids":
		ageRating = "G"
	case "teen":
		ageRating = "PG-13"
	case "adult":
		ageRating = "R"
	}

	profile := &models.Profile{
		UserID:    userID,
		Name:      name,
		Type:      profileType,
		AgeRating: ageRating,
		Avatar:    "default",
		IsMain:    isMain,
		CreatedAt: time.Now(),
	}

	err := s.profileRepo.Create(profile)
	if err != nil {
		return nil, fmt.Errorf("error al crear perfil: %w", err)
	}

	return profile, nil
}

func (s *ProfileService) GetProfilesByUserID(userID int) ([]models.Profile, error) {
	return s.profileRepo.FindByUserID(userID)
}

func (s *ProfileService) GetProfileByID(id int) (*models.Profile, error) {
	return s.profileRepo.FindByID(id)
}

func (s *ProfileService) UpdateProfile(profile *models.Profile) error {
	return s.profileRepo.Update(profile)
}

func (s *ProfileService) DeleteProfile(id int) error {
	profile, err := s.profileRepo.FindByID(id)
	if err != nil {
		return err
	}
	if profile.IsMain {
		return fmt.Errorf("no se puede eliminar el perfil principal")
	}
	return s.profileRepo.Delete(id)
}

func (s *ProfileService) CountProfilesByUserID(userID int) (int, error) {
	return s.profileRepo.CountByUserID(userID)
}