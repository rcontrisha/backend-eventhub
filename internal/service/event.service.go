package service

import (
	"context"
	"encoding/json"
	"log"
	"rcontrisha/backend-eventhub/internal/broker"
	"rcontrisha/backend-eventhub/internal/dto"
	"rcontrisha/backend-eventhub/internal/repository"
)

type EventService struct {
	repo *repository.EventRepo
	brk  *broker.Broker
}

func NewEventService(repo *repository.EventRepo, brk *broker.Broker) *EventService {
	return &EventService{
		repo: repo,
		brk:  brk,
	}
}

func (c *EventService) GetAllEvents(ctx context.Context, req dto.GetEventsRequest) (*dto.EventListResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 6
	}

	rawEvents, total, err := c.repo.GetAllEvents(ctx, req)
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

func (c *EventService) GetEventDetail(ctx context.Context, id string) (*dto.EventDetailResponse, error) {
	raw, err := c.repo.GetEventDetail(ctx, id)
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

func (c *EventService) JoinOrLeaveEvent(ctx context.Context, eventId string, userId string) (string, error) {
	joined, err := c.repo.IsJoined(ctx, eventId, userId)
	if err != nil {
		return "", err
	}

	if !joined {
		if err := c.repo.JoinEvent(ctx, eventId, userId); err != nil {
			return "", err
		}
		c.publishCount(ctx, eventId)
		return "joined", nil
	}

	if err := c.repo.LeaveEvent(ctx, eventId, userId); err != nil {
		return "", err
	}
	c.publishCount(ctx, eventId)
	return "left", nil
}

func (c *EventService) publishCount(ctx context.Context, eventId string) {
	// TODO 1: count, err := c.repo.CountParticipants(ctx, eventId)
	count, err := c.repo.CountParticipants(ctx, eventId)
	// TODO 2: kalau err -> log.Printf lalu return (jangan return error ke pemanggil)
	if err != nil {
		log.Println("[Event Publish] Failed to count event participants: ", err)
		return
	}
	// TODO 3: c.brk.Publish(dto.UpdateAttendee{EventID: eventId, Count: count})
	c.brk.Publish(dto.UpdateAttendee{EventID: eventId, Count: count})
}

func (c *EventService) GetUpcomingEvents(ctx context.Context) ([]dto.EventListItemResponse, error) {
	result, err := c.repo.UpcomingEvent(ctx)
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

func (c *EventService) GetSavedEvents(ctx context.Context, userId string) ([]dto.EventListItemResponse, error) {
	result, err := c.repo.SavedEvents(ctx, userId)
	if err != nil {
		return nil, err
	}

	savedEvents := make([]dto.EventListItemResponse, 0, len(result))
	for _, item := range result {
		var tags []string
		if err := json.Unmarshal(item.TagsRaw, &tags); err != nil {
			return nil, err
		}

		savedEvents = append(savedEvents, dto.EventListItemResponse{
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

	return savedEvents, nil
}
