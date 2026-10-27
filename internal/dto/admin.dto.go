package dto

type AdminDashboardStats struct {
	TotalUsers         int64   `json:"total_users" example:"1250"`
	UsersThisMonth     int64   `json:"users_this_month" example:"85"`
	TotalEvents        int64   `json:"total_events" example:"320"`
	UpcomingEvents     int64   `json:"upcoming_events" example:"48"`
	TotalCommunities   int64   `json:"total_communities" example:"76"`
	ActiveCommunities  int64   `json:"active_communities" example:"54"`
	AvgFillRatePercent float64 `json:"avg_fill_rate_percent" example:"68.4"`
}
