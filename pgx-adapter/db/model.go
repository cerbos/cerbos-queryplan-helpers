package db

import (
	"time"
)

type Contact struct {
	ID             int
	CreatedAt      time.Time
	UpdatedAt      time.Time
	FirstName      string
	LastName       string
	OwnerID        int
	CompanyID      int
	Active         bool
	MarketingOptIn bool
}

type Company struct {
	ID        int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	ID         int
	Username   string
	Email      string
	Name       string
	Role       string
	Department string
}
