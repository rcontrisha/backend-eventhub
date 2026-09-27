package repository

import (
	"context"
	"errors"
	"fmt"
	"rcontrisha/backend-eventhub/internal/dto"
	"rcontrisha/backend-eventhub/internal/model"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommunityRepo struct {
	db *pgxpool.Pool
}

func NewCommunityRepo(db *pgxpool.Pool) *CommunityRepo {
	return &CommunityRepo{
		db: db,
	}
}

func (c *CommunityRepo) GetAllCommunities(ctx context.Context, req dto.GetCommunityRequest) ([]model.CommunityListItem, int, error) {
	var conditions []string
	var args []any
	argIdx := 1

	if req.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(c.name ILIKE $%d OR c.description ILIKE $%d)", argIdx, argIdx))
		args = append(args, "%"+req.Search+"%")
		argIdx++
	}

	if req.Tag != "" {
		conditions = append(conditions, fmt.Sprintf(`EXISTS (
			SELECT 1 FROM community_tags ct2
			JOIN tags t2 ON t2.id = ct2.tag_id
			WHERE ct2.community_id = c.id AND t2.name ILIKE $%d
		)`, argIdx))
		args = append(args, "%"+req.Tag+"%")
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := fmt.Sprintf(`
		SELECT COUNT(DISTINCT c.id)
		FROM communities c
		%s
	`, whereClause)

	var total int
	if err := c.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.Limit
	args = append(args, req.Limit, offset)
	limitArgIdx := argIdx
	offsetArgIdx := argIdx + 1

	dataQuery := fmt.Sprintf(`
		SELECT
			c.id,
			c.name,
			c.description,
			c.banner_url,
			COALESCE(
				(
					SELECT json_agg(t.name ORDER BY t.name ASC)
					FROM community_tags ct
					JOIN tags t ON t.id = ct.tag_id
					WHERE ct.community_id = c.id
				), '[]'::json
			) AS tags,
			COALESCE(m.members_count, 0) AS members_count,
			COALESCE(e.upcoming_events, 0) AS upcoming_events,
			c.created_at,
			c.updated_at
		FROM communities c
		LEFT JOIN (
			SELECT 
				community_id, 
				COUNT(*) AS members_count
			FROM community_members
			GROUP BY community_id
		) m ON m.community_id = c.id
		LEFT JOIN (
			SELECT 
				community_id, 
				COUNT(*) AS upcoming_events
			FROM events
			WHERE start_time >= NOW()
			GROUP BY community_id
		) e ON e.community_id = c.id
		%s
		ORDER BY c.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, limitArgIdx, offsetArgIdx)

	rows, err := c.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	communities := make([]model.CommunityListItem, 0)
	for rows.Next() {
		var item model.CommunityListItem
		if err := rows.Scan(
			&item.Id,
			&item.Name,
			&item.Description,
			&item.BannerUrl,
			&item.TagsRaw,
			&item.MembersCount,
			&item.UpcomingEvents,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		communities = append(communities, item)
	}

	return communities, total, rows.Err()
}

func (c *CommunityRepo) GetCommunityDetail(ctx context.Context, communityId string) (*model.CommunityDetail, error) {
	query := `
		SELECT
			c.id,
			c.name,
			c.description,
			c.banner_url,
			COALESCE(
				(
					SELECT json_agg(t.name ORDER BY t.name ASC)
					FROM community_tags ct
					JOIN tags t ON t.id = ct.tag_id
					WHERE ct.community_id = c.id
				), '[]'::json
			) AS tags,
			COALESCE(m.members_count, 0) AS members_count,
			COALESCE(e.upcoming_events, 0) AS upcoming_events,
			c.created_at,
			c.updated_at
		FROM communities c
		LEFT JOIN (
			SELECT 
				community_id, 
				COUNT(*) AS members_count
			FROM community_members
			GROUP BY community_id
		) m ON m.community_id = c.id
		LEFT JOIN (
			SELECT 
				community_id, 
				COUNT(*) AS upcoming_events
			FROM events
			WHERE start_time >= NOW()
			GROUP BY community_id
		) e ON e.community_id = c.id
		WHERE c.id = $1
		ORDER BY c.created_at DESC
	`

	var res model.CommunityDetail

	if err := c.db.QueryRow(ctx, query, communityId).Scan(
		&res.Id,
		&res.Name,
		&res.Description,
		&res.BannerUrl,
		&res.TagsRaw,
		&res.MembersCount,
		&res.UpcomingEvents,
		&res.CreatedAt,
		&res.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("community not found")
		}
		return nil, err
	}

	return &res, nil
}

func (c *CommunityRepo) GetCommunityUpcomingEvents(ctx context.Context, communityId string) ([]model.EventListItem, error) {
	query := `
		SELECT
			e.id,
			e.title,
			e.image_url,
			COALESCE(
				(
					SELECT json_agg(t.name ORDER BY t.name ASC)
					FROM event_tags et
					JOIN tags t ON t.id = et.tag_id
					WHERE et.event_id = e.id
				), '[]'::json
			) AS tags,
			e.capacity,
			COALESCE(p.attendees_count, 0) AS attendees_count,
			e.start_time,
			e.end_time,
			e.location,
			e.created_at,
			e.updated_at
		FROM events e
		LEFT JOIN (
			SELECT 
				event_id, 
				COUNT(*) AS attendees_count
			FROM event_participants
			GROUP BY event_id
		) p ON p.event_id = e.id
		WHERE e.community_id = $1 AND e.start_time >= NOW()
		ORDER BY e.start_time ASC
	`

	rows, err := c.db.Query(ctx, query, communityId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]model.EventListItem, 0)
	for rows.Next() {
		var item model.EventListItem
		if err := rows.Scan(
			&item.Id,
			&item.Title,
			&item.ImageURL,
			&item.TagsRaw,
			&item.Capacity,
			&item.AttendeesCount,
			&item.StartTime,
			&item.EndTime,
			&item.Location,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		events = append(events, item)
	}

	return events, rows.Err()
}

func (c *CommunityRepo) GetCommunityMembers(ctx context.Context, communityId string) ([]model.CommunityMembers, error) {
	query := `
		SELECT
			u.name,
			u.avatar_url,
			cm.role
		FROM community_members cm
		JOIN communities c ON c.id = cm.community_id
		JOIN users u ON u.id = cm.user_id
		WHERE c.id = $1
	`

	rows, err := c.db.Query(ctx, query, communityId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := make([]model.CommunityMembers, 0)
	for rows.Next() {
		var member model.CommunityMembers
		if err := rows.Scan(
			&member.Name,
			&member.AvatarUrl,
			&member.CommunityRole,
		); err != nil {
			return nil, err
		}
		members = append(members, member)
	}

	return members, rows.Err()
}

func (c *CommunityRepo) GetPopularCommunities(ctx context.Context) ([]model.CommunityListItem, error) {
	query := `
		SELECT
			c.id,
			c.name,
			c.description,
			c.banner_url,
			COALESCE(
				(
					SELECT json_agg(t.name ORDER BY t.name ASC)
					FROM community_tags ct
					JOIN tags t ON t.id = ct.tag_id
					WHERE ct.community_id = c.id
				), '[]'::json
			) AS tags,
			COALESCE(m.members_count, 0) AS members_count,
			COALESCE(e.upcoming_events, 0) AS upcoming_events,
			c.created_at,
			c.updated_at
		FROM communities c
		LEFT JOIN (
			SELECT 
				community_id, 
				COUNT(*) AS members_count
			FROM community_members
			GROUP BY community_id
		) m ON m.community_id = c.id
		LEFT JOIN (
			SELECT 
				community_id, 
				COUNT(*) AS upcoming_events
			FROM events
			WHERE start_time >= NOW()
			GROUP BY community_id
		) e ON e.community_id = c.id
		ORDER BY members_count DESC
	`

	rows, err := c.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var communities []model.CommunityListItem
	for rows.Next() {
		var community model.CommunityListItem
		if err := rows.Scan(
			&community.Id,
			&community.Name,
			&community.Description,
			&community.BannerUrl,
			&community.TagsRaw,
			&community.MembersCount,
			&community.UpcomingEvents,
			&community.CreatedAt,
			&community.UpdatedAt,
		); err != nil {
			return nil, err
		}
		communities = append(communities, community)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return communities, nil
}
