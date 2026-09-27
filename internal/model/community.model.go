package model

import "time"

type Community struct {
	Id          string    `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	BannerUrl   string    `db:"banner_url"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type CommunityListItem struct {
	Id string
	Name string
	Description string
	BannerUrl string
	TagsRaw []byte
	MembersCount int
	UpcomingEvents int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CommunityMembers struct {
	CommunityId   string `db:"community_id"`
	UserId        string `db:"user_id"`
	CommunityRole string `db:"role"`
}
