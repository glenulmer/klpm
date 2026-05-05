package main

import (
	. "klpm/lib/date"
	. "klpm/lib/htmlHelper"
	. "klpm/lib/output"
)

func QuoteVarBool(v string) bool {
	switch Lower(Trim(v)) {
	case `1`, `on`, `yes`, `true`:
		return true
	}
	return false
}

func QuoteSelect(name, value string, choices []QuoteChoice_t) Elem_t {
	var options []Elem_t
	for _, x := range choices {
		opt := Option().KV(`value`, x.id).Text(x.label)
		if value == Str(x.id) { opt = opt.KV(`selected`) }
		options = append(options, opt)
	}
	return Select(options).Name(name)
}

func QuoteCheckbox(name string, checked bool) Elem_t {
	in := Elem(`input`).Type(`checkbox`).Name(name).KV(`value`, `1`)
	if checked { in = in.KV(`checked`) }
	return in
}

func QuoteInputText(name, value, placeholder string) Elem_t {
	in := Elem(`input`).Type(`text`).Name(name).Value(value).KV(`autocomplete`, `off`)
	if placeholder != `` { in = in.Place(placeholder) }
	return in
}

func QuoteInputDate(name, value string) Elem_t {
	return DateInput().
		Name(name).
		Value(value).
		KV(`inputmode`, `none`).
		KV(`onclick`, `if(this.showPicker){this.showPicker();}`).
		KV(`onfocus`, `if(this.showPicker){this.showPicker();}`).
		KV(`onkeydown`, `return false`).
		KV(`onpaste`, `return false`).
		KV(`ondrop`, `return false`)
}

func QuoteInputBuy(name, value string) Elem_t {
	minDate, maxDate, defaultDate := QuoteBuyBounds()
	buyDate := QuoteParseBuyDate(value)
	if !Valid(buyDate) { buyDate = defaultDate }
	if Valid(minDate) && int(buyDate) < int(minDate) { buyDate = minDate }
	if Valid(maxDate) && int(buyDate) > int(maxDate) { buyDate = maxDate }
	valueYMD := buyDate.Format(`yyyymmdd`)
	hidden := Elem(`input`).
		Type(`hidden`).
		Name(name).
		Value(valueYMD).
		KV(`data-date-hidden`, `1`)
	if Valid(minDate) { hidden = hidden.KV(`data-min`, minDate.Format(`yyyymmdd`)) }
	if Valid(maxDate) { hidden = hidden.KV(`data-max`, maxDate.Format(`yyyymmdd`)) }
	host := Div().
		Class(`qdate-host`).
		KV(`data-date-host`, `1`).
		KV(`value`, valueYMD)
	if Valid(minDate) { host = host.KV(`data-min`, minDate.Format(`yyyymmdd`)) }
	if Valid(maxDate) { host = host.KV(`data-max`, maxDate.Format(`yyyymmdd`)) }
	return Div().Class(`qdate`, `qdate-buy`).KV(`data-date-control`, `1`).Wrap(
		hidden,
		host,
	)
}

func QuoteInputBirth(name, value string) Elem_t {
	minDate, maxDate, defaultDate := QuoteBirthBounds()
	birthDate := QuoteParseBirthDate(value)
	if !Valid(birthDate) { birthDate = defaultDate }
	if Valid(minDate) && int(birthDate) < int(minDate) { birthDate = minDate }
	if Valid(maxDate) && int(birthDate) > int(maxDate) { birthDate = maxDate }
	valueYMD := birthDate.Format(`yyyymmdd`)
	hidden := Elem(`input`).
		Type(`hidden`).
		Name(name).
		Value(valueYMD).
		KV(`data-date-hidden`, `1`)
	if Valid(minDate) { hidden = hidden.KV(`data-min`, minDate.Format(`yyyymmdd`)) }
	if Valid(maxDate) { hidden = hidden.KV(`data-max`, maxDate.Format(`yyyymmdd`)) }
	host := Div().
		Class(`qdate-host`).
		KV(`data-date-host`, `1`).
		KV(`value`, valueYMD)
	if Valid(minDate) { host = host.KV(`data-min`, minDate.Format(`yyyymmdd`)) }
	if Valid(maxDate) { host = host.KV(`data-max`, maxDate.Format(`yyyymmdd`)) }
	return Div().Class(`qdate`, `qdate-birth`).KV(`data-date-control`, `1`).Wrap(
		hidden,
		host,
	)
}

func QuoteInputPopupDate(name string, date, minDate, maxDate CalDate_t) Elem_t {
	valueYMD := date.Format(`yyyymmdd`)
	valueText := date.Format(`dd.mm.yyyy`)
	hidden := Elem(`input`).
		Type(`hidden`).
		Name(name).
		Value(valueYMD).
		KV(`data-buy-hidden`, `1`)
	if Valid(minDate) { hidden = hidden.KV(`data-min`, minDate.Format(`yyyymmdd`)) }
	if Valid(maxDate) { hidden = hidden.KV(`data-max`, maxDate.Format(`yyyymmdd`)) }

	return Div().Class(`qbuy`).KV(`data-buy`, `1`).Wrap(
		hidden,
		Div().Class(`qbuy-trigger-wrap`).Wrap(
			Elem(`input`).
				Type(`text`).
				Class(`qbuy-trigger`).
				Value(valueText).
				KV(`readonly`).
				KV(`inputmode`, `none`).
				KV(`autocomplete`, `off`).
				KV(`data-buy-trigger`, `1`),
			Span(`▼`).Class(`qbuy-arrow`),
		),
	)
}

func QuoteInputNumber(name, value string, min, max, step int) Elem_t {
	in := Elem(`input`).Type(`number`).Name(name).Value(value)
	if min != 0 || max != 0 || step != 0 {
		in = in.KV(`min`, min).KV(`max`, max).KV(`step`, step)
	}
	return in
}

func QuoteInputSickCover(name, value, buyValue string) Elem_t {
	value = QuoteNormalizeSickCoverValue(value, buyValue)
	max := QuoteSickCoverMaxByBuyValue(buyValue)
	return Elem(`input`).
		Type(`text`).
		Name(name).
		Value(value).
		KV(`inputmode`, `numeric`).
		KV(`autocomplete`, `off`).
		KV(`data-sick-cover`, `1`).
		KV(`data-min`, 0).
		KV(`data-max`, max)
}

func QuoteControlInput(x QuoteControl_t, vars QuoteVars_t) Elem_t {
	value := QuoteValue(vars, x.name)
	switch x.kind {
	case quoteText:
		return QuoteInputText(x.name, value, x.placeholder)
	case quoteDate:
		if x.name == `birth` { return QuoteInputBirth(x.name, value) }
		if x.name == `buy` { return QuoteInputBuy(x.name, value) }
		return QuoteInputDate(x.name, value)
	case quoteNumber:
		if x.name == `sickCover` { return QuoteInputSickCover(x.name, value, QuoteValue(vars, `buy`)) }
		return QuoteInputNumber(x.name, value, x.min, x.max, x.step)
	case quoteSelect:
		return QuoteSelect(x.name, value, QuoteControlChoices(x))
	case quoteCheckbox:
		return QuoteCheckbox(x.name, QuoteVarBool(value))
	}
	return Span(`missing control`)
}

func QuoteSpanClass(span int) string {
	if span < 1 { span = 1 }
	if span > 12 { span = 12 }
	return Str(`quote-span-`, span)
}

func QuoteControlViewSpan(layout string, x QuoteControl_t, vars QuoteVars_t, span int, class ...string) Elem_t {
	if span == 0 { span = QuoteControlSpan(x, layout) }
	classes := append([]string{ QuoteSpanClass(span) }, class...)
	if x.kind == quoteCheckbox {
		return Elem(`label`).Class(`quote-check`).Class(classes...).Wrap(
			QuoteControlInput(x, vars),
			Span(x.label).Class(`quote-check-text`),
		)
	}
	return Elem(`label`).Class(`quote-field`).Class(classes...).Wrap(
		Span(x.label).Class(`quote-label`),
		QuoteControlInput(x, vars),
	)
}

func QuoteControlLabelViewSpan(layout string, x QuoteControl_t, label string, vars QuoteVars_t, span int, class ...string) Elem_t {
	x.label = label
	return QuoteControlViewSpan(layout, x, vars, span, class...)
}

func QuoteNamedControlLabelSpanView(layout, name, label string, vars QuoteVars_t, span int, class ...string) Elem_t {
	x, ok := QuoteControlByName(name)
	if !ok { return Div(`missing field: `, name) }
	return QuoteControlLabelViewSpan(layout, x, label, vars, span, class...)
}

func QuoteControlOnlyViewSpan(layout string, x QuoteControl_t, vars QuoteVars_t, span int, class ...string) Elem_t {
	if span == 0 { span = QuoteControlSpan(x, layout) }
	classes := append([]string{ QuoteSpanClass(span) }, class...)
	return Div().Class(`quote-control`).Class(classes...).Wrap(
		QuoteControlInput(x, vars),
	)
}

func QuoteNamedControlOnlySpanView(layout, name string, vars QuoteVars_t, span int, class ...string) Elem_t {
	x, ok := QuoteControlByName(name)
	if !ok { return Div(`missing field: `, name) }
	return QuoteControlOnlyViewSpan(layout, x, vars, span, class...)
}

func QuoteSpacer(span int, class ...string) Elem_t {
	classes := append([]string{ QuoteSpanClass(span), `quote-spacer` }, class...)
	return Div().Class(classes...)
}

func QuoteNamedControlSpanView(layout, name string, vars QuoteVars_t, span int, class ...string) Elem_t {
	x, ok := QuoteControlByName(name)
	if !ok { return Div(`missing field: `, name) }
	return QuoteControlViewSpan(layout, x, vars, span, class...)
}
