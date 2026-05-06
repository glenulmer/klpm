package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
	"sync"
)

const sessionCookie = `klpm_session`

type tSessionCtxKey struct{}
type tSessionCreatedCtxKey struct{}

type SessionStore_t struct {
	mu sync.RWMutex
	byToken map[string]SessionVars_t
}

var sessionCtxKey = tSessionCtxKey{}
var sessionCreatedCtxKey = tSessionCreatedCtxKey{}

func NewSessionStore() *SessionStore_t {
	return &SessionStore_t{
		byToken: make(map[string]SessionVars_t),
	}
}

func InitSessionVars() SessionVars_t {
	state := InitState()
	state.quote = QuoteDefaultVars()
	return SessionVarsFromState(state)
}

func SessionVarsFromState(state State_t) SessionVars_t {
	work := InitState()
	work.user = state.user
	work.quote = CloneQuoteVars(state.quote)
	return SessionVars_t{
		user: work.user,
		quote: work.quote,
	}
}

func StateFromSessionVars(vars SessionVars_t) State_t {
	out := InitState()
	out.quote = CloneQuoteVars(vars.quote)
	out.user = vars.user
	return out
}

func NewSessionToken() string {
	b := make([]byte, 24)
	if _, e := rand.Read(b); e != nil { panic(e) }
	return hex.EncodeToString(b)
}

func (x *SessionStore_t)EnsureToken(raw string) (token string, setCookie bool, created bool) {
	token = strings.TrimSpace(raw)

	x.mu.Lock()
	defer x.mu.Unlock()

	if token == `` {
		setCookie = true
		created = true
		for {
			token = NewSessionToken()
			if _, ok := x.byToken[token]; !ok { break }
		}
	}
	if _, ok := x.byToken[token]; !ok {
		x.byToken[token] = InitSessionVars()
		created = true
	}
	return token, setCookie, created
}

func (x *SessionStore_t)GetSessionVars(token string) SessionVars_t {
	token = strings.TrimSpace(token)
	if token == `` { return InitSessionVars() }

	x.mu.RLock()
	vars, ok := x.byToken[token]
	x.mu.RUnlock()
	if ok { return vars }
	return InitSessionVars()
}

func (x *SessionStore_t)GetState(token string) State_t {
	return StateFromSessionVars(x.GetSessionVars(token))
}

func (x *SessionStore_t)SetState(token string, state State_t) {
	token = strings.TrimSpace(token)
	if token == `` { return }
	vars := SessionVarsFromState(state)

	x.mu.Lock()
	x.byToken[token] = vars
	x.mu.Unlock()
}

func (x *SessionStore_t)MutateState(token string, fn func(*State_t)) State_t {
	token = strings.TrimSpace(token)
	state := InitState()
	if fn == nil { return state }

	x.mu.Lock()
	if token != `` {
		if vars, ok := x.byToken[token]; ok {
			state = StateFromSessionVars(vars)
		} else {
			state = StateFromSessionVars(InitSessionVars())
		}
	}
	fn(&state)
	if token != `` {
		x.byToken[token] = SessionVarsFromState(state)
	}
	x.mu.Unlock()

	return state
}

func (x *SessionStore_t)Destroy(token string) {
	token = strings.TrimSpace(token)
	if token == `` { return }
	x.mu.Lock()
	delete(x.byToken, token)
	x.mu.Unlock()
}

func SessionToken(r *http.Request) string {
	v := r.Context().Value(sessionCtxKey)
	token, ok := v.(string)
	if !ok { return `` }
	return strings.TrimSpace(token)
}

func SessionCreated(r *http.Request) bool {
	v := r.Context().Value(sessionCreatedCtxKey)
	created, ok := v.(bool)
	return ok && created
}

func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := ``
		cookie, _ := r.Cookie(sessionCookie)
		if cookie != nil { raw = cookie.Value }

		token, setCookie, created := App.sessionStore.EnsureToken(raw)
		if setCookie { SetSessionCookie(w, token) }

		ctx := context.WithValue(r.Context(), sessionCtxKey, token)
		ctx = context.WithValue(ctx, sessionCreatedCtxKey, created)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func DestroySession(r *http.Request) { App.sessionStore.Destroy(SessionToken(r)) }

func SetState(r *http.Request, state State_t) { App.sessionStore.SetState(SessionToken(r), state) }

func GetState(r *http.Request) State_t { return App.sessionStore.GetState(SessionToken(r)) }

func SetSessionCookie(w http.ResponseWriter, token string) {
	token = strings.TrimSpace(token)
	if token == `` { return }
	http.SetCookie(w, &http.Cookie{
		Name: sessionCookie,
		Value: token,
		Path: `/`,
		MaxAge: 60 * 60 * 24 * 365,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name: sessionCookie,
		Value: ``,
		Path: `/`,
		MaxAge: -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}
