package db

import (
	"time"
)

type Contact struct {
	Id             int
	CreatedAt      time.Time
	UpdatedAt      time.Time
	FirstName      string
	LastName       string
	OwnerId        int
	CompanyId      int
	Active         bool
	MarketingOptIn bool
}

type Company struct {
	Id        int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	Id         int
	Username   string
	Email      string
	Name       string
	Role       string
	Department string
}
