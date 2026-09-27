package model

import "time"

type OrganizerDashboardStats struct {
	TotalEvents    int
	TotalAttendees int
	TotalCapacity  int
}

type OrganizerDashboardEvent struct {
	ID             string
	Title          string
	ImageURL       *string
	Location       string
	StartTime      time.Time
	Capacity       int
	AttendeesCount int
}

type RegistrationChartRecord struct {
	MonthName string
	MonthDate time.Time
	Total     int
}

type OrganizerUpcomingEvent struct {
	ID             string
	Title          string
	StartTime      time.Time
	Capacity       int
	AttendeesCount int
}