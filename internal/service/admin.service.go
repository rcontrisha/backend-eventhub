package service

import (
	"context"
	"math"

	"rcontrisha/backend-eventhub/internal/dto"
	"rcontrisha/backend-eventhub/internal/repository"
)

type AdminService struct {
	repo *repository.AdminRepo
}

func NewAdminService(repo *repository.AdminRepo) *AdminService {
	return &AdminService{
		repo: repo,
	}
}

func (a *AdminService) GetDashboardStats(ctx context.Context) (*dto.AdminDashboardStats, error) {
	rawStats, err := a.repo.GetDashboardStats(ctx)
	if err != nil {
		return nil, err
	}

	var avgFillRate float64
	if rawStats.TotalCapacity > 0 {
		rate := (float64(rawStats.TotalAttendees) / float64(rawStats.TotalCapacity)) * 100
		avgFillRate = math.Round(rate*10) / 10
	}

	res := &dto.AdminDashboardStats{
		TotalUsers:         rawStats.TotalUsers,
		UsersThisMonth:     rawStats.UsersThisMonth,
		TotalEvents:        rawStats.TotalEvents,
		UpcomingEvents:     rawStats.UpcomingEvents,
		TotalCommunities:   rawStats.TotalCommunities,
		ActiveCommunities:  rawStats.ActiveCommunities,
		AvgFillRatePercent: avgFillRate,
	}

	return res, nil
}