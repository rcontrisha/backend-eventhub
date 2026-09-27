package model

import (
	"time"
)

type Event struct {
	Id          string    `db:"id"`
	Title       string    `db:"title"`
	Desc        string    `db:"desc"`
	ImageUrl    string    `db:"image_url"`
	Location    string    `db:"location"`
	StartTime   time.Time `db:"start_time"`
	EndTime     time.Time `db:"end_time"`
	Capacity    int       `db:"capacity"`
	OrganizerId string    `db:"organizer_id"`
	Communityid *string   `db:"community_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// Model to represent "Event Detail" as a result of table joins
type EventDetail struct {
	Id                string
	Title             string
	Desc              string
	ImageURL          string
	Location          string
	StartTime         time.Time
	EndTime           time.Time
	Capacity          int
	CreatedAt         time.Time
	UpdatedAt         time.Time
	TotalParticipants int
	OrganizerRaw      []byte
	CommunityRaw      []byte
	SpeakersRaw       []byte
	TagsRaw           []byte
}

type EventListItem struct {
	Id             string
	Title          string
	ImageURL       string
	TagsRaw        []byte
	Capacity       *int
	AttendeesCount int
	StartTime      time.Time
	EndTime        time.Time
	Location       string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type EventParticipant struct {
	EventID   string    `db:"event_id"`
	UserID    string    `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
