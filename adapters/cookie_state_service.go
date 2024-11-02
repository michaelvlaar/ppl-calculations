package adapters

import (
	"context"
	"encoding/json"
	"github.com/gorilla/sessions"
	"net/http"
	"os"
	"ppl-calculations/domain/state"
)

func NewCookieStateService(w http.ResponseWriter, r *http.Request) (state.Service, error) {
	return &CookieStateService{
		r:     r,
		store: sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY"))),
		w:     w,
	}, nil
}

type CookieStateService struct {
	r     *http.Request
	w     http.ResponseWriter
	store *sessions.CookieStore
	s     *state.State
}

func (service *CookieStateService) State(_ context.Context) (*state.State, error) {
	if service.s != nil {
		return service.s, nil
	}

	c, _ := service.store.Get(service.r, "_state")
	jsonState, ok := c.Values["state"]
	if ok {
		return service.newFromString(jsonState.(string))
	}

	return state.MustNew(), nil
}

func (service *CookieStateService) newFromString(jsonState string) (*state.State, error) {
	s := state.MustNew()

	err := json.Unmarshal([]byte(jsonState), &s)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (service *CookieStateService) SetState(_ context.Context, s *state.State) error {
	jsonState, err := json.Marshal(s)
	if err != nil {
		return err
	}

	service.s = s

	session, _ := service.store.Get(service.r, "_state")
	session.Values["state"] = string(jsonState)

	return service.store.Save(service.r, service.w, session)
}
