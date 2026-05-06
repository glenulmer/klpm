package main

import (
	"net/http"
	"strings"

	. "klpm/pkg.Global"
)

func RewriteQuotePage(w http.ResponseWriter, req *http.Request, state State_t) {
	QuoteEnsureDefaults(&state)
	vars := state.quote
	plans := QuotePlans(state)
	layout := QuoteRequestLayout(req)
	form := any(QuotePhoneFormBodyView(vars))
	planView := any(QuotePhonePlansView(plans))
	if layout == layoutDesktop {
		form = QuoteDesktopFormBodyView(vars)
		planView = QuoteDesktopPlansView(plans)
	}
	SendResponse(w,
		RewriteHTML(OuterHTML, `QuoteFormBody`, form),
		RewriteHTML(OuterHTML, `QuotePlans`, planView),
	)
}

func Page1QuoteChange(w http.ResponseWriter, req *http.Request) {
	name := strings.TrimSpace(req.FormValue(`name`))
	state := GetState(req)
	if name == QuoteResetControlName() {
		LoadStaticData()
		state.quote = QuoteDefaultVars()
		SetState(req, state)
		RewriteQuotePage(w, req, state)
		return
	}
	if QuoteSelectedApply(&state, name, req.FormValue(`value`)) {
		SetState(req, state)
		RewriteQuotePage(w, req, state)
		return
	}
	QuoteApply(&state, name, req.FormValue(`value`))
	SetState(req, state)
	RewriteQuotePage(w, req, state)
}
