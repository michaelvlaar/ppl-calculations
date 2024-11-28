package adapters

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"ppl-calculations/domain/state"
	"time"
)

func NewCookieStateService(w http.ResponseWriter, r *http.Request) (state.Service, error) {
	return &CookieStateService{
		r: r,
		w: w,
	}, nil
}

type CookieStateService struct {
	r *http.Request
	w http.ResponseWriter
	s *state.State
}

func (service *CookieStateService) State(_ context.Context) (*state.State, error) {
	if service.s != nil {
		return service.s, nil
	}

	if c, err := service.r.Cookie("state"); err == nil {
		return service.newFromString(c.Value)
	}

	return state.MustNew(), nil
}

func (service *CookieStateService) newFromString(base64State string) (*state.State, error) {
	s := state.MustNew()

	base64DecodedState, err := base64.StdEncoding.DecodeString(base64State)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(base64DecodedState, &s)
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

	cookie := &http.Cookie{
		Name:     "state",
		Value:    base64.StdEncoding.EncodeToString(jsonState),
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
	}

	http.SetCookie(service.w, cookie)

	return nil
}
