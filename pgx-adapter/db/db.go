// Copyright 2021-2022 Zenauth Ltd.
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"

	responsev1 "github.com/cerbos/cerbos/api/genpb/cerbos/response/v1"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
)

//go:embed seed.json
var seed []byte

type predicateBuilder interface {
	BuildPredicate(e *responsev1.ResourcesQueryPlanResponse_Expression_Operand_Expression) (where string, args []interface{}, err error)
}

type Client struct {
	client           *pgx.Conn
	predicateBuilder predicateBuilder
}

func (cli *Client) GetAllContacts(ctx context.Context) (res []*Contact, err error) {
	err = pgxscan.Select(ctx, cli.client, &res, "select * from contacts")
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (cli *Client) GetContacts(ctx context.Context, filter *responsev1.ResourcesQueryPlanResponse_Filter) (res []*Contact, err error) {
	if filter == nil {
		return nil, errors.New("\"filter\" is nil")
	}
	switch filter.Kind {
	case responsev1.ResourcesQueryPlanResponse_Filter_KIND_ALWAYS_DENIED:
		return nil, nil
	case responsev1.ResourcesQueryPlanResponse_Filter_KIND_ALWAYS_ALLOWED:
		return cli.GetAllContacts(ctx)
	case responsev1.ResourcesQueryPlanResponse_Filter_KIND_CONDITIONAL:
		where, args, err := cli.predicateBuilder.BuildPredicate(filter.Condition.GetNode().(*responsev1.ResourcesQueryPlanResponse_Expression_Operand_Expression))
		if err != nil {
			return nil, err
		}
		if where == "" {
			return cli.GetAllContacts(ctx)
		}
		err = pgxscan.Select(ctx, cli.client, &res, "select * from contacts where "+where, args...)
		if err != nil {
			return nil, err
		}
		return res, nil
	default:
		return nil, fmt.Errorf("unexpected filter.Kind: %s", filter.Kind)
	}
}

func New(ctx context.Context, b predicateBuilder, url string) (*Client, error) {
	c, err := pgx.Connect(ctx, url)
	if err != nil {
		return nil, err
	}
	return &Client{client: c, predicateBuilder: b}, nil
}

func (cli *Client) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	user := new(User)
	err := pgxscan.Get(ctx, cli.client, user, "select id, username, email, name, role, department from users where username=$1", username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

//go:embed schema.sql
var schemaSQL string

func SetupDatabase(ctx context.Context, conn *pgx.Conn) error {
	var data []Seed
	err := json.Unmarshal(seed, &data)
	if err != nil {
		return fmt.Errorf("failed to read from seed file: %w", err)
	}
	_, err = conn.Exec(ctx, schemaSQL)
	if err != nil {
		return fmt.Errorf("failed creating schema: %w", err)
	}
	for _, u := range data {
		var uid int
		err = conn.QueryRow(ctx, `
INSERT INTO USERS (name, username, email, role, department)
VALUES ($1, $2, $3, $4, $5)
RETURNING id`, u.Name, u.Username, u.Email, u.Role, u.Department).Scan(&uid)
		if err != nil {
			return err
		}
		if len(u.Contacts) > 0 {
			for _, c := range u.Contacts {
				var cid *int
				if c.Company != nil {
					var id int
					err = conn.QueryRow(ctx, `
INSERT INTO COMPANIES (name) VALUES ($1) RETURNING id`, c.Company.Name).Scan(&id)
					if err != nil {
						return fmt.Errorf("failed creating a company: %w", err)
					}
					cid = &id
				}
				_, err = conn.Exec(ctx, `
INSERT INTO CONTACTS (first_name, last_name, owner_id, company_id, active, marketing_opt_in)
VALUES ($1, $2, $3, $4, $5, $6)`, c.FirstName, c.LastName, uid, cid, c.Active, c.MarketingOptIn)

				if err != nil {
					return fmt.Errorf("failed creating a contact: %w", err)
				}
			}
		}
		if err != nil {
			return fmt.Errorf("failed creating a user: %w", err)
		}
	}
	return nil
}

type Seed struct {
	Name       string `json:"name"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	Department string `json:"department"`
	Contacts   []struct {
		Company *struct {
			Name string `json:"name"`
		} `json:"company"`
		FirstName      string `json:"firstName"`
		LastName       string `json:"lastName"`
		MarketingOptIn bool   `json:"marketingOptIn"`
		Active         bool   `json:"active"`
	} `json:"contacts,omitempty"`
}
