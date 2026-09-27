package service

import (
	"context"
	"encoding/json"
	"rcontrisha/backend-eventhub/internal/dto"
	"rcontrisha/backend-eventhub/internal/repository"
)

type EventService struct {
	repo *repository.EventRepo
}

func NewEventService(repo *repository.EventRepo) *EventService {
	return &EventService{
		repo: repo,
	}
}

func (e *EventService) GetAllEvents(ctx context.Context, req dto.GetEventsRequest) (*dto.EventListResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 10
	}

	rawEvents, total, err := e.repo.GetAllEvents(ctx, req)
	if err != nil {
		return nil, err
	}

	eventList := make([]dto.EventListItemResponse, 0, len(rawEvents))
	for _, item := range rawEvents {
		var tags []string
		if err := json.Unmarshal(item.TagsRaw, &tags); err != nil {
			return nil, err
		}

		eventList = append(eventList, dto.EventListItemResponse{
			Id:             item.Id,
			Title:          item.Title,
			ImageURL:       item.ImageURL,
			Tags:           tags,
			Capacity:       item.Capacity,
			AttendeesCount: item.AttendeesCount,
			StartTime:      item.StartTime,
			EndTime:        item.EndTime,
			Location:       item.Location,
			CreatedAt:      item.CreatedAt,
			UpdatedAt:      item.UpdatedAt,
		})
	}

	return &dto.EventListResponse{
		Events: eventList,
		Total:  total,
		Page:   req.Page,
		Limit:  req.Limit,
	}, nil
}

func (e *EventService) GetEventDetail(ctx context.Context, id string) (*dto.EventDetailResponse, error) {
	raw, err := e.repo.GetEventDetail(ctx, id)
	if err != nil {
		return nil, err
	}

	response := dto.EventDetailResponse{
		Id:                raw.Id,
		Title:             raw.Title,
		Desc:              raw.Desc,
		ImageUrl:          raw.ImageURL,
		Location:          raw.Location,
		StartTime:         raw.StartTime,
		EndTime:           raw.EndTime,
		Capacity:          raw.Capacity,
		TotalParticipants: raw.TotalParticipants,
		CreatedAt:         raw.CreatedAt,
		UpdatedAt:         raw.UpdatedAt,
	}

	if err := json.Unmarshal(raw.OrganizerRaw, &response.Organizer); err != nil {
		return nil, err
	}

	if raw.CommunityRaw != nil {
		var community dto.EventCommunityResponse
		if err := json.Unmarshal(raw.CommunityRaw, &community); err != nil {
			return nil, err
		}
		response.Community = &community
	}

	if err := json.Unmarshal(raw.SpeakersRaw, &response.Speakers); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(raw.TagsRaw, &response.Tags); err != nil {
		return nil, err
	}

	return &response, nil
}

func (e *EventService) JoinOrLeaveEvent(ctx context.Context, eventId string, userId string) (string, error) {
	joined, err := e.repo.IsJoined(ctx, eventId, userId)
	if err != nil {
		return "", err
	}

	if !joined {
		if err := e.repo.JoinEvent(ctx, eventId, userId); err != nil {
			return "", err
		}
		return "joined", nil
	}

	if err := e.repo.LeaveEvent(ctx, eventId, userId); err != nil {
		return "", err
	}
	return "left", nil
}

func (e *EventService) GetUpcomingEvents(ctx context.Context) ([]dto.EventListItemResponse, error) {
	result, err := e.repo.UpcomingEvent(ctx)
	if err != nil {
		return nil, err
	}
	eventList := make([]dto.EventListItemResponse, 0, len(result))
	for _, item := range result {
		var tags []string
		if err := json.Unmarshal(item.TagsRaw, &tags); err != nil {
			return nil, err
		}

		eventList = append(eventList, dto.EventListItemResponse{
			Id:             item.Id,
			Title:          item.Title,
			ImageURL:       item.ImageURL,
			Tags:           tags,
			Capacity:       item.Capacity,
			AttendeesCount: item.AttendeesCount,
			StartTime:      item.StartTime,
			EndTime:        item.EndTime,
			Location:       item.Location,
			CreatedAt:      item.CreatedAt,
			UpdatedAt:      item.UpdatedAt,
		})
	}

	return eventList, nil
}

func (e *EventService) GetMyEvents(ctx context.Context, userId string) ([]dto.EventListItemResponse, error) {
	result, err := e.repo.MyEvent(ctx, userId)
	if err != nil {
		return nil, err
	}

	myEvents := make([]dto.EventListItemResponse, 0, len(result))
	for _, item := range result {
		var tags []string
		if err := json.Unmarshal(item.TagsRaw, &tags); err != nil {
			return nil, err
		}

		myEvents = append(myEvents, dto.EventListItemResponse{
			Id:             item.Id,
			Title:          item.Title,
			ImageURL:       item.ImageURL,
			Tags:           tags,
			Capacity:       item.Capacity,
			AttendeesCount: item.AttendeesCount,
			StartTime:      item.StartTime,
			EndTime:        item.EndTime,
			Location:       item.Location,
			CreatedAt:      item.CreatedAt,
			UpdatedAt:      item.UpdatedAt,
		})
	}

	return myEvents, nil
}