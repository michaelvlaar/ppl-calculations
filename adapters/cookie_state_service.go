package adapters

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"github.com/gorilla/sessions"
	"net/http"
	"os"
	"ppl-calculations/domain/export"
	"ppl-calculations/domain/state"
	"slices"
	"strings"
)

const (
	CookiePrefixSecure = "__Secure-"
	CookieStateName    = "state"
	CookieExportName   = "e_"
)

func NewCookieStateService(w http.ResponseWriter, r *http.Request) (state.Service, error) {
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	store.Options.SameSite = http.SameSiteStrictMode
	store.Options.Path = "/"
	store.Options.Secure = true
	store.Options.HttpOnly = true

	cookieStateName := CookieStateName
	cookieExportName := CookieExportName

	if os.Getenv("SECURE_COOKIE") == "true" {
		cookieStateName = CookiePrefixSecure + cookieStateName
		cookieExportName = CookiePrefixSecure + cookieExportName
	}

	return &CookieStateService{
		r:                r,
		store:            store,
		w:                w,
		cookieStateName:  cookieStateName,
		cookieExportName: cookieExportName,
	}, nil
}

type CookieStateService struct {
	r                *http.Request
	w                http.ResponseWriter
	store            *sessions.CookieStore
	s                *state.State
	e                []export.Export
	cookieStateName  string
	cookieExportName string
}

func (service *CookieStateService) State(_ context.Context) (*state.State, error) {
	if service.s != nil {
		return service.s, nil
	}

	c, _ := service.store.Get(service.r, service.cookieStateName)
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

	session, _ := service.store.Get(service.r, service.cookieStateName)
	session.Values["state"] = string(jsonState)

	return service.store.Save(service.r, service.w, session)
}

func (service *CookieStateService) SetExport(ctx context.Context, e export.Export) error {
	ex, err := service.Exports(ctx)
	if err != nil {
		return err
	}

	ex = append(ex, e)

	service.e = ex

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(e); err != nil {
		return err
	}

	session, _ := service.store.Get(service.r, service.cookieExportName+e.ID.String())
	session.Values["export"] = buf.Bytes()
	session.Options.MaxAge = 60 * 60 * 24 * 30 * 6
	session.Options.HttpOnly = true
	session.Options.Secure = true
	session.Options.SameSite = http.SameSiteLaxMode
	session.Options.Path = "/"

	return service.store.Save(service.r, service.w, session)
}

func (service *CookieStateService) DeleteExport(ctx context.Context, e export.ID) error {
	ex, err := service.Exports(ctx)
	if err != nil {
		return err
	}

	service.e = ex

	index := slices.IndexFunc(service.e, func(e2 export.Export) bool {
		return e == e2.ID
	})

	if index != -1 {
		service.e = slices.Delete(service.e, index, index+1)
	}

	session, _ := service.store.Get(service.r, service.cookieExportName+e.String())
	session.Options.MaxAge = -1
	return service.store.Save(service.r, service.w, session)
}

func (service *CookieStateService) Exports(_ context.Context) ([]export.Export, error) {
	if service.e != nil {
		return service.e, nil
	}

	var exports []export.Export
	for _, c := range service.r.Cookies() {
		if strings.HasPrefix(c.Name, service.cookieExportName) {
			session, _ := service.store.Get(service.r, c.Name)
			exBytes, ok := session.Values["export"]
			if !ok {
				continue
			}

			var ex export.Export

			buf := bytes.NewBufferString(string(exBytes.([]byte)))
			dec := gob.NewDecoder(buf)
			if err := dec.Decode(&ex); err != nil {
				continue
			}

			exports = append(exports, ex)
		}
	}

	return exports, nil
}

func (service *CookieStateService) Export(_ context.Context, id export.ID) (*export.Export, error) {
	session, _ := service.store.Get(service.r, service.cookieExportName+id.String())
	exBytes, ok := session.Values["export"]
	if !ok {
		return nil, nil
	}

	var ex export.Export

	buf := bytes.NewBufferString(string(exBytes.([]byte)))
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&ex); err != nil {
		return nil, nil
	}

	return &ex, nil
}
