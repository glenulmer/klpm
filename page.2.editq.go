package main

import (
	"net/http"

	. "klpm/lib/htmlHelper"
	. "klpm/lib/output"
)

func EditQCSSPath() string {
	return Str(`/static/css/review.css?v=`, App.staticVersion)
}

func EditQBodyView(vars QuoteVars_t, sortForGet bool) Elem_t {
	return Div().Id(`EditQFormBody`).Class(`editq-body`).Wrap(
		EditQHeaderView(vars),
		EditQDependantsView(vars, sortForGet),
		EditQQuoteReviewCardView(vars),
	)
}

func EditQFormView(vars QuoteVars_t, sortForGet bool) Elem_t {
	return Elem(`form`).
		Id(`EditQForm`).
		Class(`editq-form`).
		KV(`method`, `post`).
		KV(`action`, `/quote-review-change`).
		Wrap(EditQBodyView(vars, sortForGet))
}

func EditQPageView(vars QuoteVars_t, sortForGet bool) Elem_t {
	return Elem(`main`).Class(`editq-page`).Wrap(
		EditQFormView(vars, sortForGet),
	)
}

func Page2EditQ(w0 http.ResponseWriter, req *http.Request) {
	state := GetState(req)
	QuoteEnsureDefaults(&state)
	if len(QuoteSelectedItems(state.quote)) == 0 { http.Redirect(w0, req, `/`, http.StatusSeeOther); return }
	EditQDropPreinsertedDependant(&state)
	SetState(req, state)

	head := Head()
	head = head.
		CSS(Str(`/static/css/common.css?v=`, App.staticVersion)).
		CSS(EditQCSSPath()).
		CSSTail(Str(`/static/css/date-ctrl.css?v=`, App.staticVersion)).
		JSTail(Str(`/static/js/date-ctrl.js?v=`, App.staticVersion)).
		JSTail(Str(`/static/js/page.2.editq.js?v=`, App.staticVersion)).
		Title(SiteName).
		End()

	w := Writer(w0)
	w.Add(
		head.Left(), NL,
		EditQPageView(state.quote, true), NL,
		head.Right(), NL,
	)
}
