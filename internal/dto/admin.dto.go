package dto

type AdminDashboardStats struct {
	TotalUsers         int64   `json:"total_users"`
	UsersThisMonth     int64   `json:"users_this_month"`
	TotalEvents        int64   `json:"total_events"`
	UpcomingEvents     int64   `json:"upcoming_events"`
	TotalCommunities   int64   `json:"total_communities"`
	ActiveCommunities  int64   `json:"active_communities"`
	AvgFillRatePercent float64 `json:"avg_fill_rate_percent"`
}
