package service

import (
	"net/http"
	"github.com/cerbos/cerbos-go-adapters/ent-adapter/ent"
	"log"
	"context"
	cerbos "github.com/cerbos/cerbos/client"
	"encoding/json"
	responsev1 "github.com/cerbos/cerbos/api/genpb/cerbos/response/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

type (
	authCtxKeyType struct{}

	authContext struct {
		username  string
		principal *cerbos.Principal
	}

	dbClient interface {
		GetUserByUsername(ctx context.Context, username string) (*ent.User, error)
		GetContacts(ctx context.Context, filter *responsev1.ResourcesQueryPlanResponse_Filter) ([]*ent.Contact, error)
	}

	Service struct {
		cerbos cerbos.Client
		client dbClient
	}
)

var authCtxKey = authCtxKeyType{}

func (svc *Service) getContacts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authCtx := ctx.Value(authCtxKey).(*authContext)
	result, err := svc.cerbos.ResourcesQueryPlan(ctx, authCtx.principal, cerbos.NewResource("contact", ""), "read")
	if err != nil {
		writeMessage(w, http.StatusInternalServerError, err.Error())
		return
	}
	filter := result.GetFilter()
	log.Println(protojson.Format(filter))
	res, err := svc.client.GetContacts(ctx, filter)
	if err != nil {
		writeMessage(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (svc *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	svc.authenticationMiddleware(http.HandlerFunc(svc.getContacts)).ServeHTTP(w, r)
}

func NewService(cerbosAddr string, client dbClient) (*Service, error) {
	c, err := cerbos.New(cerbosAddr, cerbos.WithPlaintext())
	if err != nil {
		return nil, err
	}
	return &Service{
		cerbos: c,
		client: client,
	}, nil
}

// authenticationMiddleware handles the verification of username and password,
// creates a Cerbos principal and adds it to the request context.
func (svc *Service) authenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the basic auth credentials from the request.
		user, _, ok := r.BasicAuth()
		if ok {
			// retrieve the auth context.
			authCtx, err := svc.buildAuthContext(user, r)
			if err != nil {
				log.Printf("Failed to authenticate user [%s]: %v", user, err)
			} else {
				// Add the retrieved principal to the context.
				ctx := context.WithValue(r.Context(), authCtxKey, authCtx)
				next.ServeHTTP(w, r.WithContext(ctx))

				return
			}
		}

		// No credentials provided or the credentials are invalid.
		w.Header().Set("WWW-Authenticate", `Basic realm="auth"`)
		writeMessage(w, http.StatusUnauthorized, "Authentication required")
	})
}

// buildAuthContext verifies the username and password and returns a new authContext object.
func (svc *Service) buildAuthContext(username string, r *http.Request) (*authContext, error) {
	// Lookup the user from the database.
	user, err := svc.client.GetUserByUsername(r.Context(), username)
	if err != nil {
		return nil, err
	}

	// Create a new principal object with information from the database and the request.
	principal := cerbos.NewPrincipal(username).
		WithRoles(user.Role).
		WithAttr("department", user.Department).
		WithAttr("ipAddress", r.RemoteAddr)

	return &authContext{username: username, principal: principal}, nil
}

type genericResponse struct {
	Message string `json:"message"`
}

func writeMessage(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, genericResponse{Message: msg})
}

func writeJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")

	_ = enc.Encode(v)
}
