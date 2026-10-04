package repository

import (
	"context"
	"fmt"
	"log"
	"rcontrisha/backend-eventhub/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminRepo struct {
	db *pgxpool.Pool
}

func NewAdminRepo(db *pgxpool.Pool) *AdminRepo {
	return &AdminRepo{
		db: db,
	}
}

func (a *AdminRepo) GetDashboardStats(ctx context.Context) (*model.AdminDashboardStats, error) {
	// belom get actual 'active communities'
	query := `
		SELECT 
			(SELECT COUNT(*) FROM users) AS total_users,
			(SELECT COUNT(*) FROM users WHERE created_at >= DATE_TRUNC('month', CURRENT_TIMESTAMP)) AS users_this_month,
			(SELECT COUNT(*) FROM events) AS total_events,
			(SELECT COUNT(*) FROM events WHERE start_time > CURRENT_TIMESTAMP) AS upcoming_events,
			(SELECT COUNT(*) FROM communities) AS total_communities,
			(SELECT COUNT(*) FROM communities) AS active_communities, 
			(SELECT COALESCE(SUM(capacity), 0) FROM events) AS total_capacity,
			(SELECT COUNT(*) FROM event_participants) AS total_attendees;
	`

	var stats model.AdminDashboardStats
	err := a.db.QueryRow(ctx, query).Scan(
		&stats.TotalUsers,
		&stats.UsersThisMonth,
		&stats.TotalEvents,
		&stats.UpcomingEvents,
		&stats.TotalCommunities,
		&stats.ActiveCommunities,
		&stats.TotalCapacity,
		&stats.TotalAttendees,
	)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to fetch admin dashboard stats: %w", err)
	}

	return &stats, nil
}
