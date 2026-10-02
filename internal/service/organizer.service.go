package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"math"
	"rcontrisha/backend-eventhub/internal/dto"
	"rcontrisha/backend-eventhub/internal/model"
	"rcontrisha/backend-eventhub/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OrganizerService struct {
	repo *repository.OrganizerRepo
	db   *pgxpool.Pool
}

func NewOrganizerService(repo *repository.OrganizerRepo, db *pgxpool.Pool) *OrganizerService {
	return &OrganizerService{
		repo: repo,
		db:   db,
	}
}

func (s *OrganizerService) GetDashboard(ctx context.Context, organizerId string) (*dto.OrganizerDashboardResponse, error) {
	rawStats, err := s.repo.GetDashboardStats(ctx, s.db, organizerId)
	if err != nil {
		return nil, err
	}

	rawEvents, err := s.repo.GetYourEvents(ctx, s.db, organizerId)
	if err != nil {
		return nil, err
	}

	rawChart, err := s.repo.GetRegistrationChart(ctx, s.db, organizerId)
	if err != nil {
		return nil, err
	}

	rawUpcoming, err := s.repo.GetUpcomingEventsMini(ctx, s.db, organizerId)
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

func (o *OrganizerService) AddEvent(ctx context.Context, organizerId string, payload dto.AddEventRequest) error {
	tx, err := o.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			log.Println(err.Error())
		}
	}()

	var speakersJson []model.Speaker
	jsonErr := json.Unmarshal([]byte(payload.Speakers), &speakersJson)
	log.Println(jsonErr)
	if jsonErr != nil {
		return errors.New("invalid json format for speakers")
	}

	var data = model.Event{
		Title:       payload.Title,
		Desc:        payload.Desc,
		ImageUrl:    payload.ImageUrl,
		CommunityId: payload.CommunityId,
		StartTime:   payload.StartTime,
		EndTime:     payload.EndTime,
		Location:    payload.Location,
		Capacity:    payload.Capacity,
		Speakers:    speakersJson,
	}

	eventId, err := o.repo.InsertEvent(ctx, tx, organizerId, data)
	if err != nil {
		return err
	}

	if _, err := o.repo.InsertEventTags(ctx, tx, eventId, payload.Tags); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
