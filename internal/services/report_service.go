package services

import (
	"SDGEStreaming/internal/db"
	"SDGEStreaming/internal/repositories"
	"fmt"
)

type ReportService struct {
	userRepo            repositories.UserRepo
	contentRepo         repositories.ContentRepo
	playbackHistoryRepo repositories.PlaybackHistoryRepo
	subscriptionRepo    repositories.SubscriptionRepo
}

type UserReport struct {
	TotalUsers        int
	ActiveUsers       int
	UsersByPlan       map[string]int
	NewUsersThisMonth int
}

type ContentReport struct {
	TotalAudiovisual  int
	TotalAudio        int
	TopRatedContent   []string
	MostViewedContent []string
}

type RevenueReport struct {
	TotalRevenue  float64
	RevenueByPlan map[string]float64
	Transactions  int
}

func NewReportService(
	userRepo repositories.UserRepo,
	contentRepo repositories.ContentRepo,
	playbackHistoryRepo repositories.PlaybackHistoryRepo,
	subscriptionRepo repositories.SubscriptionRepo,
) *ReportService {
	return &ReportService{
		userRepo:            userRepo,
		contentRepo:         contentRepo,
		playbackHistoryRepo: playbackHistoryRepo,
		subscriptionRepo:    subscriptionRepo,
	}
}

func (s *ReportService) GenerateUserReport() (*UserReport, error) {
	users, err := s.userRepo.FindAll()
	if err != nil {
		return nil, err
	}

	report := &UserReport{
		TotalUsers:  len(users),
		UsersByPlan: make(map[string]int),
	}

	planNames := map[int]string{1: "Free", 2: "Estandar", 3: "Premium 4K"}
	for _, user := range users {
		planName := planNames[user.PlanID]
		report.UsersByPlan[planName]++
	}

	return report, nil
}

func (s *ReportService) GenerateContentReport() (*ContentReport, error) {
	audiovisuals, _ := s.contentRepo.FindAllAudiovisual()
	audios, _ := s.contentRepo.FindAllAudio()

	report := &ContentReport{
		TotalAudiovisual:  len(audiovisuals),
		TotalAudio:        len(audios),
		TopRatedContent:   []string{},
		MostViewedContent: []string{},
	}

	for i, av := range audiovisuals {
		if i >= 5 {
			break
		}
		if av.AverageRating > 0 {
			report.TopRatedContent = append(report.TopRatedContent,
				fmt.Sprintf("%s (%.1f/10)", av.Title, av.AverageRating))
		}
	}

	return report, nil
}

func (s *ReportService) GenerateRevenueReport() (*RevenueReport, error) {
	conn := db.GetDB()

	var totalRevenue float64
	err := conn.QueryRow("SELECT COALESCE(SUM(amount), 0) FROM payment_transactions WHERE status = 'completed'").Scan(&totalRevenue)
	if err != nil {
		totalRevenue = 0
	}

	var transactions int
	conn.QueryRow("SELECT COUNT(*) FROM payment_transactions").Scan(&transactions)

	report := &RevenueReport{
		TotalRevenue:  totalRevenue,
		RevenueByPlan: make(map[string]float64),
		Transactions:  transactions,
	}

	return report, nil
}