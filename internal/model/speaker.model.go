package model

type Speaker struct {
	Id   string `db:"id"`
	Name string `db:"name"`
	Role string `db:"role"`
}