package model

type AdminDashboardStats struct {
	TotalUsers        int64
	UsersThisMonth    int64
	TotalEvents       int64
	UpcomingEvents    int64
	TotalCommunities  int64
	ActiveCommunities int64
	TotalCapacity     int64
	TotalAttendees    int64
}