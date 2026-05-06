package main

import (
	"net/http"
	"strings"

	. "klpm/pkg.Global"
)

func RewriteEditQPage(w http.ResponseWriter, req *http.Request, state State_t) {
	QuoteEnsureDefaults(&state)
	vars := state.quote
	body := EditQBodyView(vars, false)
	SendResponse(w, RewriteHTML(OuterHTML, `EditQFormBody`, body))
}

func Page2EditQEntry(w http.ResponseWriter, req *http.Request) {
	state := GetState(req)
	QuoteEnsureDefaults(&state)
	EditQDropPreinsertedDependant(&state)
	SetState(req, state)
	http.Redirect(w, req, `/quote-review`, http.StatusSeeOther)
}

func Page2EditQChange(w http.ResponseWriter, req *http.Request) {
	name := strings.TrimSpace(req.FormValue(`name`))
	value := req.FormValue(`value`)
	token := SessionToken(req)
	applied := false
	state := App.sessionStore.MutateState(token, func(state *State_t) {
		applied = EditQApply(state, name, value)
	})
	if applied {
		RewriteEditQPage(w, req, state)
		return
	}

	http.Redirect(w, req, `/quote-review`, http.StatusSeeOther)
}
