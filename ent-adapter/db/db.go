package db

import (
	"github.com/cerbos/cerbos-go-adapters/ent-adapter/ent"
	_ "embed"
	_ "github.com/mattn/go-sqlite3"
	"context"
	"encoding/json"
	"fmt"
	"entgo.io/ent/dialect"
	"github.com/cerbos/cerbos-go-adapters/ent-adapter/ent/user"
	responsev1 "github.com/cerbos/cerbos/api/genpb/cerbos/response/v1"
	"errors"
	"entgo.io/ent/dialect/sql"
)

//go:embed seed.json
var seed []byte

type predicateBuilder interface {
	BuildPredicate(e *responsev1.ResourcesQueryPlanResponse_Expression_Operand_Expression) (p *sql.Predicate, err error)
}

type Client struct {
	client           *ent.Client
	predicateBuilder predicateBuilder
}

func (cli *Client) GetContacts(ctx context.Context, filter *responsev1.ResourcesQueryPlanResponse_Filter) ([]*ent.Contact, error) {
	if filter == nil {
		return nil, errors.New("\"filter\" is nil")
	}
	switch filter.Kind {
	case responsev1.ResourcesQueryPlanResponse_Filter_KIND_ALWAYS_DENIED:
		return nil, nil
	case responsev1.ResourcesQueryPlanResponse_Filter_KIND_ALWAYS_ALLOWED:
		return cli.client.Contact.Query().All(ctx)
	case responsev1.ResourcesQueryPlanResponse_Filter_KIND_CONDITIONAL:
		p, err := cli.predicateBuilder.BuildPredicate(filter.Condition.GetNode().(*responsev1.ResourcesQueryPlanResponse_Expression_Operand_Expression))
		if err != nil {
			return nil, err
		}
		if p == nil {
			return cli.client.Contact.Query().All(ctx)
		}
		return cli.client.Contact.Query().Where(func(s *sql.Selector) {
			s.Where(p)
		}).All(ctx)
	}
	return nil, errors.New("unknown filter kind")
}

func New(b predicateBuilder, options ...ent.Option) (*Client, error) {
	c, err := ent.Open(dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1", options...)
	if err != nil {
		return nil, err
	}
	return &Client{client: c, predicateBuilder: b}, nil
}

func (cli *Client) GetUserByUsername(ctx context.Context, username string) (*ent.User, error) {
	user, err := cli.client.User.Query().Where(user.UsernameEQ(username)).Only(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (cli *Client) SetupDatabase(ctx context.Context) error {
	var data []Seed
	err := json.Unmarshal(seed, &data)
	if err != nil {
		return fmt.Errorf("failed to read from seed file: %w", err)
	}
	// Run the automatic migration tool to create all schema resources.
	client := cli.client
	if err := client.Schema.Create(ctx); err != nil {
		return fmt.Errorf("failed creating schema resources: %w", err)
	}
	for _, u := range data {
		user, err := client.User.Create().
			SetName(u.Name).
			SetUsername(u.Username).
			SetEmail(u.Email).
			SetRole(u.Role).
			SetDepartment(u.Department).
			Save(ctx)

		if len(u.Contacts) > 0 {
			for _, c := range u.Contacts {
				var cid *int
				if c.Company != nil {
					c, err := client.Company.Create().
						SetName(c.Company.Name).
						Save(ctx)
					if err != nil {
						return fmt.Errorf("failed creating a company: %w", err)
					}
					cid = &c.ID
				}
				_, err := client.Contact.Create().
					SetFirstName(c.FirstName).
					SetLastName(c.LastName).
					SetMarketingOptIn(c.MarketingOptIn).
					SetActive(c.Active).
					SetOwner(user).
					SetNillableCompanyID(cid).
					Save(ctx)

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
		FirstName      string `json:"firstName"`
		LastName       string `json:"lastName"`
		MarketingOptIn bool   `json:"marketingOptIn"`
		Active         bool   `json:"active"`
		Company        *struct {
			Name string `json:"name"`
		} `json:"company"`
	} `json:"contacts,omitempty"`
}
