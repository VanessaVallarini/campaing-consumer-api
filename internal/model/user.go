package model

import "time"

type UserDBModel struct {
	Id        string
	Email     string
	CreatedAt *time.Time
	UpdatedAt *time.Time
	Active    bool
}
