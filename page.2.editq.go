package main

import (
	"net/http"

	. "klpm/lib/htmlHelper"
	. "klpm/lib/output"
)

func EditQCSSPath(layout string) string {
	if layout == layoutDesktop {
		return Str(`/static/css/page.2.editq.desktop.css?v=`, App.staticVersion)
	}
	return Str(`/static/css/page.2.editq.phone.css?v=`, App.staticVersion)
}

func EditQBodyView(layout string, vars QuoteVars_t, sortForGet bool) Elem_t {
	if layout == layoutDesktop { return EditQDesktopBodyView(vars, sortForGet) }
	return EditQPhoneBodyView(vars, sortForGet)
}

func EditQFormView(layout string, vars QuoteVars_t, sortForGet bool) Elem_t {
	return Elem(`form`).
		Id(`EditQForm`).
		Class(`editq-form`).
		KV(`method`, `post`).
		KV(`action`, `/quote-review-change`).
		Wrap(EditQBodyView(layout, vars, sortForGet))
}

func EditQPageView(layout string, vars QuoteVars_t, sortForGet bool) Elem_t {
	return Elem(`main`).Class(`editq-page`).Wrap(
		EditQFormView(layout, vars, sortForGet),
	)
}

func Page2EditQ(w0 http.ResponseWriter, req *http.Request) {
	state := GetState(req)
	QuoteEnsureDefaults(&state)
	if len(QuoteSelectedItems(state.quote)) == 0 { http.Redirect(w0, req, `/`, http.StatusSeeOther); return }
	EditQDropPreinsertedDependant(&state)
	SetState(req, state)
	layout := RequestLayout(req)

	head := Head()
	head = head.
		CSS(EditQCSSPath(layout)).
		JSTail(Str(`/static/js/page.2.editq.js?v=`, App.staticVersion)).
		Title(SiteName).
		End()

	w := Writer(w0)
	w.Add(
		head.Left(), NL,
		EditQPageView(layout, state.quote, true), NL,
		head.Right(), NL,
	)
}
