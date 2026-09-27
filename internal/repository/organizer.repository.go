package repository

import (
	"context"
	"fmt"
	"rcontrisha/backend-eventhub/internal/model"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OrganizerRepo struct {
	db *pgxpool.Pool
}

func NewOrganizerRepo(db *pgxpool.Pool) *OrganizerRepo {
	return &OrganizerRepo{
		db: db,
	}
}

func (o *OrganizerRepo) GetDashboardStats(ctx context.Context, organizerId string) (*model.OrganizerDashboardStats, error) {
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
	err := o.db.QueryRow(ctx, query, organizerId).Scan(
		&stats.TotalEvents,
		&stats.TotalAttendees,
		&stats.TotalCapacity,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stats: %w", err)
	}

	return &stats, nil
}

func (r *OrganizerRepo) GetYourEvents(ctx context.Context, organizerId string) ([]model.OrganizerDashboardEvent, error) {
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

	rows, err := r.db.Query(ctx, query, organizerId)
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

func (r *OrganizerRepo) GetRegistrationChart(ctx context.Context, organizerID string) ([]model.RegistrationChartRecord, error) {
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

	rows, err := r.db.Query(ctx, query, organizerID)
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

func (r *OrganizerRepo) GetUpcomingEventsMini(ctx context.Context, organizerId string) ([]model.OrganizerUpcomingEvent, error) {
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

	rows, err := r.db.Query(ctx, query, organizerId)
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
