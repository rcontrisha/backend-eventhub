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
