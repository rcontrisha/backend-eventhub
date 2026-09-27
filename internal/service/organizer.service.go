package service

import (
	"context"
	"math"
	"rcontrisha/backend-eventhub/internal/dto"
	"rcontrisha/backend-eventhub/internal/repository"
)

type OrganizerService struct {
	repo *repository.OrganizerRepo
}

func NewOrganizerService(repo *repository.OrganizerRepo) *OrganizerService {
	return &OrganizerService{
		repo: repo,
	}
}

func (s *OrganizerService) GetDashboard(ctx context.Context, organizerId string) (*dto.OrganizerDashboardResponse, error) {
	rawStats, err := s.repo.GetDashboardStats(ctx, organizerId)
	if err != nil {
		return nil, err
	}

	rawEvents, err := s.repo.GetYourEvents(ctx, organizerId)
	if err != nil {
		return nil, err
	}

	rawChart, err := s.repo.GetRegistrationChart(ctx, organizerId)
	if err != nil {
		return nil, err
	}

	rawUpcoming, err := s.repo.GetUpcomingEventsMini(ctx, organizerId)
	if err != nil {
		return nil, err
	}

	var avgFillRate float64
	if rawStats.TotalCapacity > 0 {
		rate := (float64(rawStats.TotalAttendees) / float64(rawStats.TotalCapacity)) * 100
		avgFillRate = math.Round(rate*10) / 10
	}

	statsDTO := dto.DashboardStatsResponse{
		TotalEvents:    rawStats.TotalEvents,
		TotalAttendees: rawStats.TotalAttendees,
		AvgFillRate:    avgFillRate,
		EventViews:     3241, 
	}

	yourEventsDTO := make([]dto.DashboardEventItemResponse, 0, len(rawEvents))
	for _, item := range rawEvents {
		yourEventsDTO = append(yourEventsDTO, dto.DashboardEventItemResponse{
			ID:             item.ID,
			Title:          item.Title,
			ImageURL:       item.ImageURL,
			Location:       item.Location,
			StartTime:      item.StartTime,
			Capacity:       item.Capacity,
			AttendeesCount: item.AttendeesCount,
		})
	}

	chartDTO := make([]dto.DashboardChartItem, 0, len(rawChart))
	for _, item := range rawChart {
		chartDTO = append(chartDTO, dto.DashboardChartItem{
			Month: item.MonthName,
			Total: item.Total,
		})
	}

	upcomingDTO := make([]dto.DashboardUpcomingItemResponse, 0, len(rawUpcoming))
	for _, item := range rawUpcoming {
		upcomingDTO = append(upcomingDTO, dto.DashboardUpcomingItemResponse{
			ID:             item.ID,
			Title:          item.Title,
			StartTime:      item.StartTime,
			Capacity:       item.Capacity,
			AttendeesCount: item.AttendeesCount,
		})
	}

	return &dto.OrganizerDashboardResponse{
		Stats:          statsDTO,
		ChartData:      chartDTO,
		YourEvents:     yourEventsDTO,
		UpcomingEvents: upcomingDTO,
	}, nil
}