package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"rcontrisha/backend-eventhub/internal/dto"
	"rcontrisha/backend-eventhub/internal/model"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EventRepo struct {
	db *pgxpool.Pool
}

func NewEventRepo(db *pgxpool.Pool) *EventRepo {
	return &EventRepo{
		db: db,
	}
}

func (e *EventRepo) GetAllEvents(ctx context.Context, req dto.GetEventsRequest) ([]model.EventListItem, int, error) {
	var conditions []string
	var args []any
	argIdx := 1

	// Filter pencarian title
	if req.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(e.title ILIKE $%d)", argIdx))
		args = append(args, "%"+req.Search+"%")
		argIdx++
	}

	// Filter lokasi
	if req.Location != "" {
		conditions = append(conditions, fmt.Sprintf("e.location ILIKE $%d", argIdx))
		args = append(args, "%"+req.Location+"%")
		argIdx++
	}

	// Filter tag (memastikan event memiliki tag tertentu)
	if req.Tag != "" {
		conditions = append(conditions, fmt.Sprintf(`EXISTS (
			SELECT 1 FROM event_tags et2
			JOIN tags t2 ON t2.id = et2.tag_id
			WHERE et2.event_id = e.id AND t2.name ILIKE $%d
		)`, argIdx))
		args = append(args, "%"+req.Tag+"%")
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// 1. Query Total Count untuk pagination
	countQuery := fmt.Sprintf(`
		SELECT COUNT(DISTINCT e.id)
		FROM events e
		%s
	`, whereClause)

	var total int
	if err := e.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// 2. Query Data dengan pagination
	offset := (req.Page - 1) * req.Limit
	args = append(args, req.Limit, offset)
	limitArgIdx := argIdx
	offsetArgIdx := argIdx + 1

	dataQuery := fmt.Sprintf(`
		SELECT
			e.id,
			e.title,
			e.image_url,
			COALESCE(
				(
					SELECT json_agg(t.name)
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
			SELECT event_id, COUNT(*) AS attendees_count
			FROM event_participants
			GROUP BY event_id
		) p ON p.event_id = e.id
		%s
		ORDER BY e.start_time ASC
		LIMIT $%d OFFSET $%d
	`, whereClause, limitArgIdx, offsetArgIdx)

	rows, err := e.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	events := make([]model.EventListItem, 0)
	for rows.Next() {
		var item model.EventListItem
		if err := rows.Scan(
			&item.ID,
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
			return nil, 0, err
		}
		events = append(events, item)
	}

	return events, total, rows.Err()
}

func (e *EventRepo) GetEventDetail(ctx context.Context, id string) (*model.EventDetail, error) {
	query := `
		SELECT 
			e.id,
			e.title,
			e.desc,
			e.image_url,
			e.location,
			e.start_time,
			e.end_time,
			e.capacity,
			e.created_at,
			e.updated_at,
			COALESCE(p.total_participants, 0) AS total_participants,
			json_build_object(
				'id', u.id,
				'name', u.name,
				'email', u.email,
				'avatar_url', u.avatar_url,
				'role', u.role
			) AS organizer,
			CASE 
				WHEN c.id IS NOT NULL THEN json_build_object(
					'id', c.id,
					'name', c.name,
					'description', c.description,
					'banner_url', c.banner_url
				)
				ELSE NULL 
			END AS community,
			COALESCE(
				(
					SELECT json_agg(json_build_object('id', s.id, 'name', s.name, 'role', s.role))
					FROM event_speakers es
					JOIN speakers s ON es.speaker_id = s.id
					WHERE es.event_id = e.id
				), '[]'::json
			) AS speakers,
			COALESCE(
				(
					SELECT json_agg(json_build_object('id', t.id, 'name', t.name))
					FROM event_tags et
					JOIN tags t ON et.tag_id = t.id
					WHERE et.event_id = e.id
				), '[]'::json
			) AS tags
		FROM events e
		JOIN users u ON e.organizer_id = u.id
		LEFT JOIN communities c ON e.community_id = c.id
		LEFT JOIN (
			SELECT event_id, COUNT(*) AS total_participants
			FROM event_participants
			GROUP BY event_id
		) p ON p.event_id = e.id
		WHERE e.id = $1
	`
	var res model.EventDetail

	if err := e.db.QueryRow(ctx, query, id).Scan(
		&res.Id,
		&res.Title,
		&res.Desc,
		&res.ImageURL,
		&res.Location,
		&res.StartTime,
		&res.EndTime,
		&res.Capacity,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.TotalParticipants,
		&res.OrganizerRaw,
		&res.CommunityRaw,
		&res.SpeakersRaw,
		&res.TagsRaw,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("event not found")
		}
		return nil, err
	}

	return &res, nil
}

func (e *EventRepo) IsJoined(ctx context.Context, eventId string, userId string) (bool, error) {
	log.Println(userId)
	log.Println(eventId)
	query := `SELECT EXISTS(SELECT 1 FROM event_participants WHERE event_id = $1 AND user_id = $2)`
	var exists bool
	if err := e.db.QueryRow(ctx, query, eventId, userId).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (e *EventRepo) JoinEvent(ctx context.Context, eventId string, userId string) error {
	query := `
		INSERT INTO event_participants (event_id, user_id) 
		VALUES ($1, $2)
	`
	cmd, err := e.db.Exec(ctx, query, eventId, userId)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("user already joined or event not found")
	}
	return nil
}

func (e *EventRepo) LeaveEvent(ctx context.Context, eventId string, userId string) error {
	query := `DELETE FROM event_participants WHERE event_id = $1 AND user_id = $2`
	cmd, err := e.db.Exec(ctx, query, eventId, userId)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("user is not participating in this event")
	}
	return nil
}
