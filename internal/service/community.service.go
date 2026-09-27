package service

import (
	"context"
	"encoding/json"
	"rcontrisha/backend-eventhub/internal/dto"
	"rcontrisha/backend-eventhub/internal/repository"
)

type CommunityService struct {
	repo *repository.CommunityRepo
}

func NewCommunityService(repo *repository.CommunityRepo) *CommunityService {
	return &CommunityService{
		repo: repo,
	}
}

func (c *CommunityService) GetAllCommunities(ctx context.Context, req dto.GetCommunityRequest) (*dto.CommunityListResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 10
	}

	rawCommunities, total, err := c.repo.GetAllCommunities(ctx, req)
	if err != nil {
		return nil, err
	}

	communityList := make([]dto.CommunityListItemResponse, 0, len(rawCommunities))
	for _, item := range rawCommunities {
		var tags []string
		if len(item.TagsRaw) > 0 {
			if err := json.Unmarshal(item.TagsRaw, &tags); err != nil {
				return nil, err
			}
		} else {
			tags = []string{}
		}

		communityList = append(communityList, dto.CommunityListItemResponse{
			Id:             item.Id,
			Name:           item.Name,
			Description:    item.Description,
			BannerUrl:      item.BannerUrl,
			Tags:           tags,
			MembersCount:   item.MembersCount,
			UpcomingEvents: item.UpcomingEvents,
			CreatedAt:      item.CreatedAt,
			UpdatedAt:      item.UpdatedAt,
		})
	}

	return &dto.CommunityListResponse{
		Communities: communityList,
		Total:       total,
		Page:        req.Page,
		Limit:       req.Limit,
	}, nil
}

func (c *CommunityService) GetCommunityDetail(ctx context.Context, communityId string) (*dto.CommunityDetailResponse, error) {
	raw, err := c.repo.GetCommunityDetail(ctx, communityId)
	if err != nil {
		return nil, err
	}

	response := dto.CommunityDetailResponse{
		Id:             raw.Id,
		Name:           raw.Name,
		Description:    raw.Description,
		BannerUrl:      raw.BannerUrl,
		MembersCount:   raw.MembersCount,
		UpcomingEvents: raw.UpcomingEvents,
		CreatedAt:      raw.CreatedAt,
		UpdatedAt:      raw.UpdatedAt,
	}

	if err := json.Unmarshal(raw.TagsRaw, &response.Tags); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *CommunityService) GetCommunityUpcomingEvents(ctx context.Context, communityId string) ([]dto.EventListItemResponse, error) {
	result, err := c.repo.GetCommunityUpcomingEvents(ctx, communityId)
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
