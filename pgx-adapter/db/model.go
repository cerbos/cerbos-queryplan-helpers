// Copyright 2021-2022 Zenauth Ltd.
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"time"
)

type Contact struct {
	CreatedAt      time.Time
	UpdatedAt      time.Time
	FirstName      string
	LastName       string
	ID             int
	OwnerID        int
	CompanyID      int
	Active         bool
	MarketingOptIn bool
}

type Company struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	ID        int
}

type User struct {
	Username   string
	Email      string
	Name       string
	Role       string
	Department string
	ID         int
}
