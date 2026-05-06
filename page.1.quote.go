package main

import (
	"net/http"

	. "klpm/lib/date"
	. "klpm/lib/htmlHelper"
	. "klpm/lib/output"
)

func QuoteCSSPath(layout string) string {
	if layout == layoutDesktop {
		return Str(`/static/css/page.1.quote.desktop.css?v=`, App.staticVersion)
	}
	return Str(`/static/css/page.1.quote.phone.css?v=`, App.staticVersion)
}

func QuotePageView(layout string, vars QuoteVars_t, plans QuotePlans_t) Elem_t {
	if layout == layoutDesktop { return QuoteDesktopPageView(vars, plans) }
	return QuotePhonePageView(vars, plans)
}

func queryYesNo(value string) bool {
	switch Left(Lower(Trim(value)), 1) {
	case `n`, `0`, ``:
		return false
	}
	return true
}

func querySegment(value string) string {
	switch Lower(Left(Trim(value), 4)) {
	case `stud`, `4`:
		return `4`
	case `free`, `2`:
		return `2`
	}
	return `1`
}

func Page1QuoteApplyQuery(state *State_t, req *http.Request) {
	q := req.URL.Query()
	if len(q) == 0 { return }
	QuoteEnsureDefaults(state)

	set := func(name, value string) { QuoteSetValue(&state.quote, name, value) }
	setBool := func(name string, on bool) {
		if on { set(name, `1`) } else { set(name, ``) }
	}
	swapRange := func(minKey, maxKey string) {
		v1, ok1 := StateIntOK(*state, minKey)
		v2, ok2 := StateIntOK(*state, maxKey)
		if !ok1 || !ok2 || v1 <= v2 { return }
		set(minKey, Str(v2))
		set(maxKey, Str(v1))
	}

	const levelBand = 10
	const CH, CD = catHospital*levelBand, catDental*levelBand
	fn, ln := ``, ``

	for key, values := range q {
		value := ``
		if len(values) > 0 { value = Trim(values[0]) }

		switch key {
		case `dob`:
			raw := Atoi(value)
			dob := DateFromYMD(raw/10000, (raw/100)%100, raw%100)
			latest := CurrentDBDate()
			earliest := DateFromYMD(latest.Year()-75, 1, 1)
			if int(dob) > int(latest) { dob = latest }
			if int(dob) < int(earliest) { dob = earliest }
			set(`birth`, dob.Format(`yyyymmdd`))

		case `cover`:
			set(`sickCover`, Str(OnlyDigits(value)))

		case `employed`:
			set(`segment`, `1`)
		case `segment`:
			set(`segment`, querySegment(value))

		case `examOK`:
			if queryYesNo(value) { set(`exam`, `0`) } else { set(`exam`, `1`) }
		case `glasses`:
			setBool(`vision`, queryYesNo(value))
		case `natural`:
			setBool(`naturalMed`, queryYesNo(value))
		case `visa`:
			setBool(`tempVisa`, value == `limited`)

		case `doctor`:
			set(`specref`, Str(2-Atoi(value)))
		case `referral`:
			set(`specref`, Str(Atoi(value)))

		case `firstName`:
			fn = value
		case `lastName`:
			ln = value
		case `email`:
			set(`email`, value)

		case `minDeduct`:
			set(`deductibleMin`, Str(OnlyDigits(value)))
		case `maxDeduct`:
			set(`deductibleMax`, Str(OnlyDigits(value)))

		case `minHospital`:
			set(`hospitalMin`, Str(CH + Atoi(value)%levelBand))
		case `maxHospital`:
			set(`hospitalMax`, Str(CH + Atoi(value)%levelBand))

		case `minDental`:
			set(`dentalMin`, Str(CD + Atoi(value)%levelBand))
		case `maxDental`:
			set(`dentalMax`, Str(CD + Atoi(value)%levelBand))
		}
	}
	if fn != `` || ln != `` {
		set(`clientName`, Trim(Str(fn, ` `, ln)))
	}

	swapRange(`deductibleMin`, `deductibleMax`)
	swapRange(`hospitalMin`, `hospitalMax`)
	swapRange(`dentalMin`, `dentalMax`)
}

func Page1Quote(w0 http.ResponseWriter, req *http.Request) {
	state := GetState(req)
	QuoteEnsureDefaults(&state)
	Page1QuoteApplyQuery(&state, req)
	SetState(req, state)
	plans := QuotePlans(state)
	layout := RequestLayout(req)

	head := Head()
	head = head.
		CSS(QuoteCSSPath(layout)).
		CSSTail(Str(`/static/css/page.1.quote.date.css?v=`, App.staticVersion)).
		JSTail(Str(`/static/js/page.1.quote.date.js?v=`, App.staticVersion)).
		JSTail(Str(`/static/js/page.1.quote.js?v=`, App.staticVersion)).
		Title(SiteName).
		End()

	w := Writer(w0)
	w.Add(
		head.Left(), NL,
		QuotePageView(layout, state.quote, plans), NL,
		head.Right(), NL,
	)
}
