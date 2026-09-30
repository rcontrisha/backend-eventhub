package repository

import (
	"context"
	"fmt"
	"log"
	"rcontrisha/backend-eventhub/internal/dto"
	"rcontrisha/backend-eventhub/internal/model"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DBTX interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type OrganizerRepo struct{}

func NewOrganizerRepo() *OrganizerRepo {
	return &OrganizerRepo{}
}

func (o *OrganizerRepo) GetDashboardStats(ctx context.Context, db DBTX, organizerId string) (*model.OrganizerDashboardStats, error) {
	query := `
		SELECT 
			COUNT(e.id) AS total_events,
			COALESCE(SUM(p.attendees_count), 0) AS total_attendees,
			COALESCE(SUM(e.capacity), 0) AS total_capacity
		FROM events e
		LEFT JOIN (
			SELECT event_id, COUNT(*) AS attendees_count
			FROM event_participants
			GROUP BY event_id
		) p ON p.event_id = e.id
		WHERE e.organizer_id = $1
	`

	var stats model.OrganizerDashboardStats
	err := db.QueryRow(ctx, query, organizerId).Scan(
		&stats.TotalEvents,
		&stats.TotalAttendees,
		&stats.TotalCapacity,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stats: %w", err)
	}

	return &stats, nil
}

func (o *OrganizerRepo) GetYourEvents(ctx context.Context, db DBTX, organizerId string) ([]model.OrganizerDashboardEvent, error) {
	query := `
		SELECT 
			e.id, 
			e.title, 
			e.image_url, 
			e.location, 
			e.start_time, 
			e.capacity, 
			COALESCE(p.attendees_count, 0) AS attendees_count
		FROM events e
		LEFT JOIN (
			SELECT event_id, COUNT(*) AS attendees_count
			FROM event_participants
			GROUP BY event_id
		) p ON p.event_id = e.id
		WHERE e.organizer_id = $1
		ORDER BY e.start_time DESC
	`

	rows, err := db.Query(ctx, query, organizerId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch your events: %w", err)
	}
	defer rows.Close()

	events := make([]model.OrganizerDashboardEvent, 0)
	for rows.Next() {
		var item model.OrganizerDashboardEvent
		if err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.ImageURL,
			&item.Location,
			&item.StartTime,
			&item.Capacity,
			&item.AttendeesCount,
		); err != nil {
			return nil, fmt.Errorf("failed to scan your event row: %w", err)
		}
		events = append(events, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row error on your events: %w", err)
	}

	return events, nil
}

func (o *OrganizerRepo) GetRegistrationChart(ctx context.Context, db DBTX, organizerID string) ([]model.RegistrationChartRecord, error) {
	query := `
		SELECT 
			TO_CHAR(DATE_TRUNC('month', ep.created_at), 'Mon') AS month_name,
			DATE_TRUNC('month', ep.created_at) AS month_date,
			COUNT(ep.user_id) AS total
		FROM event_participants ep
		JOIN events e ON e.id = ep.event_id
		WHERE e.organizer_id = $1 
		  AND ep.created_at >= DATE_TRUNC('month', NOW() - INTERVAL '5 months')
		GROUP BY month_date, month_name
		ORDER BY month_date ASC
	`

	rows, err := db.Query(ctx, query, organizerID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch registration chart: %w", err)
	}
	defer rows.Close()

	records := make([]model.RegistrationChartRecord, 0)
	for rows.Next() {
		var item model.RegistrationChartRecord
		var monthDate time.Time
		if err := rows.Scan(
			&item.MonthName,
			&monthDate,
			&item.Total,
		); err != nil {
			return nil, fmt.Errorf("failed to scan chart record: %w", err)
		}
		item.MonthDate = monthDate
		records = append(records, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row error on chart records: %w", err)
	}

	return records, nil
}

func (o *OrganizerRepo) GetUpcomingEventsMini(ctx context.Context, db DBTX, organizerId string) ([]model.OrganizerUpcomingEvent, error) {
	query := `
		SELECT 
			e.id, 
			e.title, 
			e.start_time, 
			e.capacity, 
			COALESCE(p.attendees_count, 0) AS attendees_count
		FROM events e
		LEFT JOIN (
			SELECT event_id, COUNT(*) AS attendees_count
			FROM event_participants
			GROUP BY event_id
		) p ON p.event_id = e.id
		WHERE e.organizer_id = $1 AND e.start_time >= NOW()
		ORDER BY e.start_time ASC
	`

	rows, err := db.Query(ctx, query, organizerId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch mini upcoming events: %w", err)
	}
	defer rows.Close()

	upcoming := make([]model.OrganizerUpcomingEvent, 0)
	for rows.Next() {
		var item model.OrganizerUpcomingEvent
		if err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.StartTime,
			&item.Capacity,
			&item.AttendeesCount,
		); err != nil {
			return nil, fmt.Errorf("failed to scan upcoming event row: %w", err)
		}
		upcoming = append(upcoming, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row error on upcoming events: %w", err)
	}

	return upcoming, nil
}

func (o *OrganizerRepo) InsertEvent(ctx context.Context, db DBTX, organizerId string, payload dto.AddEventDataRequest) (string, error) {
	log.Println(payload)
	query := `INSERT INTO events ("title", "desc", "image_url", "location", "start_time", "end_time", "organizer_id", "community_id", "capacity", "speakers") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`
	args := []any{payload.Title, payload.Desc, payload.ImageUrl, payload.Location, payload.StartTime, payload.EndTime, organizerId, payload.CommunityId, payload.Capacity, payload.Speakers}

	var id string
	if err := db.QueryRow(ctx, query, args...).Scan(&id); err != nil {
		log.Println(err.Error())
		return "", err
	}

	log.Printf("New Event's id: %s", id)
	return id, nil
}

func (o *OrganizerRepo) InsertEventTags(ctx context.Context, db DBTX, eventId string, tagIds dto.AddEventTagsRequest) (pgconn.CommandTag, error) {
	query := "INSERT INTO event_tags (event_id, tag_id) VALUES "
	args := []any{}
	for idx, tagId := range tagIds.Tags {
		num := (idx * 2) + 1
		query += fmt.Sprintf("($%d, $%d)", num, num+1)
		args = append(args, eventId, tagId)
		if idx < len(tagIds.Tags)-1 {
			query += ","
		}
	}

	return db.Exec(ctx, query, args...)
}
