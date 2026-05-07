package main

import (
	. "klpm/lib/dec2"
	. "klpm/lib/htmlHelper"
	. "klpm/lib/output"
)

func EditQDelButton(name string) Elem_t {
	return Elem(`button`).
		Type(`button`).
		Name(name).
		Value(`1`).
		Class(`editq-del-btn`).
		Text(`🗑`)
}

func EditQAddButton(name, label string, class ...string) Elem_t {
	return Elem(`button`).
		Type(`button`).
		Name(name).
		Value(`1`).
		Class(`editq-add-btn`).
		Class(class...).
		Text(label)
}

func EditQPreexChargeModeView(name, mode string) Elem_t {
	return Select(
		Option().KV(`value`, editQPreexModePct).Text(`%`),
		Option().KV(`value`, editQPreexModeEur).Text(`€`),
	).Name(name).Choose(EditQPreexMode(mode)).Class(`editq-preex-mode`)
}

func EditQPreexChargeRowView(x EditQPreexCharge_t, modeKey, amountKey, noteKey string) Elem_t {
	return Div().Class(`editq-preex-row`).Wrap(
		Div(x.level).Class(`editq-preex-categ`),
		Div().Class(`editq-preex-inputs`).Wrap(
			EditQPreexChargeModeView(modeKey, x.mode),
			QuoteInputText(amountKey, x.amount, `Amount`).Class(`editq-preex-amount`, `editq-preex-amount-input`),
			QuoteInputText(noteKey, x.note, `Optional note`).Class(`editq-preex-note`),
		),
	)
}

func EditQEuroCentText(v EuroCent_t) string {
	return Str(v.String(), ` €`)
}

func EditQEuroCentOrDash(v EuroCent_t) string {
	return v.OutEuro()
}

func EditQPreexTotal(vars QuoteVars_t) EuroCent_t {
	charges := EditQPreexCharges(vars)
	var total EuroCent_t
	for _, x := range charges { total += x.applied }
	return total
}

func EditQPreexTitleText(vars QuoteVars_t) string {
	totalFlat := EuroFlatFromCent(EditQPreexTotal(vars))
	return Str(`Pre-existing conditions (`, totalFlat.OutEuro(), `)`)
}

func EditQPreexChargesView(vars QuoteVars_t) Elem_t {
	charges := EditQPreexCharges(vars)
	appliedByItem := make(map[int]EuroCent_t)
	planByItem := make(map[int]EuroCent_t)
	for _, x := range charges {
		appliedByItem[x.itemId] += x.applied
		if _, ok := planByItem[x.itemId]; !ok { planByItem[x.itemId] = x.planPrice }
	}

	var rows []Elem_t
	prevItemId := 0
	hasPlan := false
	for _, x := range charges {
		if x.itemId != prevItemId {
			if hasPlan {
				base := planByItem[prevItemId]
				preex := appliedByItem[prevItemId]
				rows = append(rows, Div().Class(`editq-preex-row`, `editq-preex-summary`).Wrap(
					Div(EditQEuroCentText(base)).Class(`editq-preex-summary-base`),
					Div().Class(`editq-preex-inputs`, `editq-preex-summary-inputs`).Wrap(
						Div().Class(`editq-preex-summary-spacer`),
						Div(EditQEuroCentText(preex)).Class(`editq-preex-summary-preex`),
						Div(Str(`= `, EditQEuroCentText(base+preex))).Class(`editq-preex-summary-sum`),
					),
				))
			}
			rows = append(rows, Div().Class(`editq-preex-plan`).Wrap(
				Span(x.plan).Class(`editq-preex-plan-name`),
			))
			prevItemId = x.itemId
			hasPlan = true
		}
		rows = append(rows, EditQPreexChargeRowView(x, EditQPreexModeKey(x.itemId, x.categId), EditQPreexAmountKey(x.itemId, x.categId), EditQPreexNoteKey(x.itemId, x.categId)))
	}
	if hasPlan {
		base := planByItem[prevItemId]
		preex := appliedByItem[prevItemId]
		rows = append(rows, Div().Class(`editq-preex-row`, `editq-preex-summary`).Wrap(
			Div(EditQEuroCentText(base)).Class(`editq-preex-summary-base`),
			Div().Class(`editq-preex-inputs`, `editq-preex-summary-inputs`).Wrap(
				Div().Class(`editq-preex-summary-spacer`),
				Div(EditQEuroCentText(preex)).Class(`editq-preex-summary-preex`),
				Div(Str(`= `, EditQEuroCentText(base+preex))).Class(`editq-preex-summary-sum`),
			),
		))
	}
	if len(rows) == 0 {
		rows = append(rows, Div(`No payable selected plan categories yet.`).Class(`editq-preex-empty`))
	}
	return Div().Class(`editq-section`, `editq-preex-charges`).Wrap(
		Div(`Pre-existing conditions charges`).Class(`editq-section-title`),
		Div().Class(`editq-preex-list`).Wrap(rows),
	)
}

func EditQDependantView(vars QuoteVars_t, dep EditQDep_t) Elem_t {
	charges := EditQDependantCharges(vars, dep)
	appliedByItem := make(map[int]EuroCent_t)
	planByItem := make(map[int]EuroCent_t)
	for _, x := range charges {
		appliedByItem[x.itemId] += x.applied
		if _, ok := planByItem[x.itemId]; !ok { planByItem[x.itemId] = x.planPrice }
	}

	var depRows []Elem_t
	prevItemId := 0
	hasPlan := false
	for _, x := range charges {
		if x.itemId != prevItemId {
			if hasPlan {
				base := planByItem[prevItemId]
				preex := appliedByItem[prevItemId]
				depRows = append(depRows, Div().Class(`editq-preex-row`, `editq-preex-summary`).Wrap(
					Div(EditQEuroCentText(base)).Class(`editq-preex-summary-base`),
					Div().Class(`editq-preex-inputs`, `editq-preex-summary-inputs`).Wrap(
						Div().Class(`editq-preex-summary-spacer`),
						Div(EditQEuroCentText(preex)).Class(`editq-preex-summary-preex`),
						Div(Str(`= `, EditQEuroCentText(base+preex))).Class(`editq-preex-summary-sum`),
					),
				))
			}
			ageText := ``
			if x.planAge > 0 {
				ageText = Str(`Effective age: `, x.planAge)
				if x.planAgeMode == `exact` { ageText = Str(`Effective age (Exact): `, x.planAge) }
				if x.planAgeMode == `year` { ageText = Str(`Effective age (Year): `, x.planAge) }
			}
			if ageText == `` {
				depRows = append(depRows, Div().Class(`editq-preex-plan`).Wrap(
					Span(x.plan).Class(`editq-preex-plan-name`),
				))
			} else {
				depRows = append(depRows, Div().Class(`editq-preex-plan`).Wrap(
					Div().Class(`editq-preex-plan-head`).Wrap(
						Span(x.plan).Class(`editq-preex-plan-name`),
					),
					Div(ageText).Class(`editq-preex-plan-age`),
				))
			}
			prevItemId = x.itemId
			hasPlan = true
		}
		depRows = append(depRows, EditQPreexChargeRowView(
			x,
			EditQDepChargeModeKey(dep.depId, x.itemId, x.categId),
			EditQDepChargeAmountKey(dep.depId, x.itemId, x.categId),
			EditQDepChargeNoteKey(dep.depId, x.itemId, x.categId),
		))
	}
	if hasPlan {
		base := planByItem[prevItemId]
		preex := appliedByItem[prevItemId]
		depRows = append(depRows, Div().Class(`editq-preex-row`, `editq-preex-summary`).Wrap(
			Div(EditQEuroCentText(base)).Class(`editq-preex-summary-base`),
			Div().Class(`editq-preex-inputs`, `editq-preex-summary-inputs`).Wrap(
				Div().Class(`editq-preex-summary-spacer`),
				Div(EditQEuroCentText(preex)).Class(`editq-preex-summary-preex`),
				Div(Str(`= `, EditQEuroCentText(base+preex))).Class(`editq-preex-summary-sum`),
			),
		))
	}
	if len(depRows) == 0 {
		depRows = append(depRows, Div(`No selected plan available.`).Class(`editq-preex-empty`))
	}
	return Div().Class(`editq-dependant`).Wrap(
		Div().Class(`editq-dependant-fields-row`).Wrap(
			QuoteInputText(EditQDepNameKey(dep.depId), dep.name, `Dependant name`).Class(`editq-dep-name`),
			QuoteInputBirth(EditQDepBirthKey(dep.depId), dep.birth).Class(`editq-dep-birth`),
			Elem(`label`).Class(`editq-check`, `editq-dep-vision`).KV(`title`, `Vision`).Wrap(
				QuoteCheckbox(EditQDepVisionKey(dep.depId), dep.vision),
				Span(`Vision`).Class(`editq-check-text`),
			),
			EditQDelButton(EditQDepDelControlName(dep.depId)).Class(`editq-dep-del`),
		),
		Div().Class(`editq-preex-charges`, `editq-dependant-preex`).Wrap(
			Div(`Pre-existing conditions charges`).Class(`editq-section-title`),
			Div().Class(`editq-preex-list`).Wrap(depRows),
		),
	)
}

func EditQDependantsView(vars QuoteVars_t, sortForGet bool) Elem_t {
	deps := EditQDependants(vars, sortForGet)
	depCount := len(deps)
	var list []Elem_t
	for _, dep := range deps {
		list = append(list, EditQDependantView(vars, dep))
	}
	if len(deps) < editQDepMaxCount {
		list = append(list, EditQAddButton(EditQDepAddControlName(), `Add Dependant`, `editq-add-dependant`))
	}
	return EditQTopCardView(
		`EditQDependantsCard`,
		Str(`Dependants (`, depCount, `)`),
		false,
		Div().Class(`editq-section`, `editq-dependants`).Wrap(
			Div().Class(`editq-dependant-list`).Wrap(list),
		),
	)
}

func EditQTopCardView(id, title string, open bool, body Elem_t, right ...Elem_t) Elem_t {
	var titleRow []Elem_t
	titleRow = append(titleRow, Span(title).Class(`editq-card-title-text`))
	if len(right) > 0 {
		titleRow = append(titleRow, Div().Class(`editq-card-title-right`).Wrap(right))
	}
	card := Elem(`details`).Id(id).Class(`editq-fold-card`)
	if open { card = card.KV(`open`, `open`) }
	return card.Wrap(
		Elem(`summary`).Class(`editq-card-title`, `editq-fold-title`).Wrap(
			Div().Class(`editq-card-title-row`).Wrap(titleRow),
		),
		Div().Class(`editq-card-body`).Wrap(body),
	)
}

func EditQReviewControlValue(vars QuoteVars_t, name string) string {
	value := QuoteValue(vars, name)
	if value == `` { return `-` }
	ctrl, ok := QuoteControlByName(name)
	if !ok { return value }
	if ctrl.kind == quoteCheckbox {
		if QuoteVarBool(value) { return `Yes` }
		return `No`
	}
	if ctrl.kind != quoteSelect { return value }
	for _, x := range QuoteControlChoices(ctrl) {
		if Str(x.id) == value { return x.label }
	}
	return value
}

func EditQReviewYear(value string) string {
	value = Trim(value)
	if len(value) < 4 { return `-` }
	return value[:4]
}

func EditQReviewClientView(vars QuoteVars_t) Elem_t {
	clientName := Trim(QuoteValue(vars, `clientName`))
	if clientName == `` { clientName = `No client name` }
	return Div().Class(`editq-review-client`).Wrap(
		Div(clientName).Class(`editq-review-client-name`),
		Div().Class(`editq-review-client-grid`).Wrap(
			Div(`Year Born`).Class(`editq-review-client-key`),
			Div(EditQReviewYear(QuoteValue(vars, `birth`))).Class(`editq-review-client-val`),
			Div(`Year Bought`).Class(`editq-review-client-key`),
			Div(EditQReviewYear(QuoteValue(vars, `buy`))).Class(`editq-review-client-val`),
			Div(`Sick Cover`).Class(`editq-review-client-key`),
			Div(EditQReviewControlValue(vars, `sickCover`)).Class(`editq-review-client-val`),
			Div(`Client Type`).Class(`editq-review-client-key`),
			Div(EditQReviewControlValue(vars, `segment`)).Class(`editq-review-client-val`),
			Div(`Vision`).Class(`editq-review-client-key`),
			Div(EditQReviewControlValue(vars, `vision`)).Class(`editq-review-client-val`),
		),
	)
}

func EditQReviewPreexByItemCateg(vars QuoteVars_t) (map[int]EuroCent_t, map[string]EuroCent_t) {
	byItem := make(map[int]EuroCent_t)
	byItemCateg := make(map[string]EuroCent_t)
	for _, x := range EditQPreexCharges(vars) {
		key := Str(x.itemId, `:`, x.categId)
		byItem[x.itemId] += x.applied
		byItemCateg[key] += x.applied
	}
	return byItem, byItemCateg
}

func EditQReviewPriceRow(label string, total, base, surch, preex EuroCent_t, class ...string) Elem_t {
	return Div().Class(`editq-review-row`).Class(class...).Wrap(
		Div(label).Class(`editq-review-col-label`),
		Div(EditQEuroCentOrDash(total)).Class(`editq-review-col-money`),
		Div(EditQEuroCentOrDash(base)).Class(`editq-review-col-money`),
		Div(EditQEuroCentOrDash(surch)).Class(`editq-review-col-money`),
		Div(EditQEuroCentOrDash(preex)).Class(`editq-review-col-money`, `editq-review-col-preex`),
	)
}

func EditQReviewPlanView(vars QuoteVars_t, item QuoteSelectedItem_t, row QuotePlan_t, preexByItem map[int]EuroCent_t, preexByItemCateg map[string]EuroCent_t) Elem_t {
	preexTotal := preexByItem[item.itemId]
	preexPlan := preexByItemCateg[Str(item.itemId, `:`, 0)]
	effectiveAge := `-`
	ageLabel := `Age`
	work := QuoteStateFromQuoteVars(vars)
	_, yearAge, exactAge := PlanAges(work)
	age := yearAge
	if EditQPlanUsesExactAge(row.planId) {
		age = exactAge
		ageLabel = `Exact age`
	}
	if age > 0 { effectiveAge = Str(age) }
	var rows []Elem_t
	rows = append(rows, Div().Class(`editq-review-row`, `editq-review-row-head`).Wrap(
		Div(`Item`).Class(`editq-review-col-label`),
		Div(`Total`).Class(`editq-review-col-money`),
		Div(`Base`).Class(`editq-review-col-money`),
		Div(`Surch`).Class(`editq-review-col-money`),
		Div(`Pre-ex`).Class(`editq-review-col-money`, `editq-review-col-preex`),
	))
	rows = append(rows, EditQReviewPriceRow(`Total`, row.base+row.surcharge+preexTotal, row.base, row.surcharge, preexTotal, `editq-review-row-total`))
	rows = append(rows, EditQReviewPriceRow(`Plan`, row.planBase+row.planSurcharge+preexPlan, row.planBase, row.planSurcharge, preexPlan))
	for _, addon := range row.addons {
		if !addon.priceOk && addon.addon == 0 && addon.level == 0 && addon.label == `` { continue }
		label := QuoteAddonPickText(addon)
		if label == `` { label = addon.categ }
		preex := preexByItemCateg[Str(item.itemId, `:`, addon.categId)]
		rows = append(rows, EditQReviewPriceRow(label, addon.base+addon.surcharge+preex, addon.base, addon.surcharge, preex))
	}

	return Div().Class(`editq-review-plan`).Wrap(
		Div(row.label).Class(`editq-review-plan-title`),
		Div().Class(`editq-review-plan-meta`).Wrap(
			Span(Str(ageLabel, ` `, effectiveAge)),
			Span(Str(`Ded `, PriceTextWholeEuro(row.deduct, true))),
			Span(Str(`NC `, PriceTextWholeEuro(row.noClaims, true))),
			Span(Str(`Comm `, PriceTextWholeEuro(row.commission, true))),
		),
		Div().Class(`editq-review-grid`).Wrap(rows),
	)
}

func EditQReviewExportButtonsView(vars QuoteVars_t) Elem_t {
	state := QuoteStateFromQuoteVars(vars)
	quote := QuoteVars(&state)

	langs := []Elem_t{}
	langDefault := Str(int(quote.lang))
	slimDefault := `0`
	if quote.slim { slimDefault = `1` }
	for id, x := range App.lookup.languages.All() {
		if len(langs) == 0 && langDefault == `0` { langDefault = Str(id) }
		langs = append(langs, Option().KV(`value`, id).Text(x.label))
	}
	if len(langs) == 0 { langs = append(langs, Option().KV(`value`, `0`).Text(`Default`)) }

	return Div().
		Class(`editq-card-title-right`).
		KV(`style`, `width: 100%; justify-content: space-between;`).
		Wrap(
			Div().KV(`style`, `display: flex; align-items: center; gap: 0.3rem; flex-wrap: wrap;`).Wrap(
				Select(langs).Name(`lang`).Choose(langDefault),
				Select(
					Option().KV(`value`, `0`).Text(`Complete`),
					Option().KV(`value`, `1`).Text(`Slim`),
				).Name(`slim`).Choose(slimDefault),
			),
			Elem(`button`).
				Type(`button`).
				Name(`DownloadExcel`).
				Value(`1`).
				Class(`editq-title-btn`).
				Text(`Get Excel`),
		)
}

func EditQQuoteReviewBody(vars QuoteVars_t) Elem_t {
	state := QuoteStateFromQuoteVars(vars)
	selected := QuoteSelectedRows(state)
	preexByItem, preexByItemCateg := EditQReviewPreexByItemCateg(vars)

	var plans []Elem_t
	for _, x := range selected {
		plans = append(plans, EditQReviewPlanView(vars, x.item, x.row, preexByItem, preexByItemCateg))
	}
	if len(plans) == 0 {
		plans = append(plans, Div(`No plans selected.`).Class(`editq-review-empty`))
	}
	return Div().Class(`editq-review`).Wrap(
		EditQReviewExportButtonsView(vars),
		EditQReviewClientView(vars),
		Div().Class(`editq-review-plans`).Wrap(plans),
	)
}

func EditQHeaderView(vars QuoteVars_t) Elem_t {
	return EditQTopCardView(`EditQPreexCard`, EditQPreexTitleText(vars), false,
		Div().Class(`editq-header`).Wrap(
			EditQPreexChargesView(vars),
		),
	)
}

func EditQQuoteReviewCardView(vars QuoteVars_t) Elem_t {
	backBtn := Elem(`a`).
		KV(`href`, `/`).
		Class(`editq-title-btn`).
		Text(`Back to Quote Info`)
	return EditQTopCardView(`EditQReviewCard`, `Quote Review`, true,
		EditQQuoteReviewBody(vars),
		backBtn,
	)
}
