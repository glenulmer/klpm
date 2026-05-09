package main

import (
	"fmt"
	"strings"
	. "klpm/lib/htmlHelper"
	. "klpm/lib/dec2"
	. "klpm/lib/output"
)

func QuotePlanAddonByTag(x QuotePlan_t, tag string) (QuotePlanAddon_t, bool) {
	tag = Lower(Trim(tag))
	for _, addon := range x.addons {
		if Contains(Lower(Trim(addon.categ)), tag) { return addon, true }
	}
	return QuotePlanAddon_t{}, false
}

func QuotePlanAddonByCateg(x QuotePlan_t, categId CategId_t) (QuotePlanAddon_t, bool) {
	for _, addon := range x.addons {
		if addon.categId == categId { return addon, true }
	}
	return QuotePlanAddon_t{}, false
}

func QuoteAddonPickText(x QuotePlanAddon_t) string {
	if x.label != `` && Lower(Trim(x.label)) != Lower(Trim(x.categ)) { return x.label }
	if x.addon != 0 { return AddonName(CatChoice_t{ addon:x.addon, label:x.label }) }
	if x.level != 0 { return LevelLabel(x.level) }
	return ``
}

func QuotePlanAddonSelectView(planId int, x QuotePlanAddon_t) Elem_t {
	return QuotePlanAddonSelectNamedView(QuotePlanCatControlName(planId, x.categId), x)
}

func QuotePlanAddonSelectNamedView(name string, x QuotePlanAddon_t) Elem_t {
	var options []Elem_t
	for _, choice := range x.choices {
		options = append(options, Option().KV(`value`, choice.addon).Text(AddonName(choice)))
	}
	return Select(options).Id(name).Name(name).Choose(x.addon)
}

func QuotePlanAddonView(planId int, x QuotePlanAddon_t) Elem_t {
	pick := QuoteAddonPickText(x)
	pickView := any(pick)
	if x.hasMulti && len(x.choices) > 0 {
		pickView = QuotePlanAddonSelectView(planId, x)
	}
	return Div().Class(`quote-plan-addon`).Wrap(
		Div(x.categ).Class(`quote-plan-addon-cat`),
		Div().Class(`quote-plan-addon-pick`).Wrap(pickView),
		Div(PriceText(x.base, x.priceOk)).Class(`quote-plan-addon-base`),
		Div(PriceText(x.surcharge, x.priceOk)).Class(`quote-plan-addon-surch`),
	)
}

func QuotePlanBaseView(x QuotePlan_t) Elem_t {
	return Div().Class(`quote-plan-addon`, `quote-plan-addon-base-row`).Wrap(
		Div(`Plan`).Class(`quote-plan-addon-cat`),
		Div(``).Class(`quote-plan-addon-pick`),
		Div(PriceText(x.planBase, x.planOk)).Class(`quote-plan-addon-base`),
		Div(PriceText(x.planSurcharge, x.planOk)).Class(`quote-plan-addon-surch`),
	)
}

func QuotePlanCardView(x QuotePlan_t) Elem_t {
	var addons []Elem_t
	addons = append(addons, QuotePlanBaseView(x))
	for _, addon := range x.addons { addons = append(addons, QuotePlanAddonView(x.planId, addon)) }
	return Div().Class(`quote-plan-card`).KV(`data-plan-id`, Str(x.planId)).Wrap(
		Div().Class(`quote-plan-head`).Wrap(
			Div().Class(`quote-plan-head-main`).Wrap(
				Elem(`button`).Type(`button`).Name(QuoteSelectedAddControlName(x.planId)).Value(Str(x.planId)).Class(`quote-plan-pick-btn`, `quote-plan-pick-add`).Text(`🛒`),
				Div(x.label).Class(`quote-plan-label`),
			),
			Div().Class(`quote-plan-total`).Wrap(
				Div(PriceTextWholeEuro(x.price, true)).Class(`quote-plan-price`),
			),
		),
		Div().Class(`quote-plan-meta`).Wrap(
			Div().Class(`quote-plan-meta-item`).Wrap(
				Span(`Ded`).Class(`quote-plan-meta-label`),
				Span(PriceTextWholeEuro(x.deduct, true)).Class(`quote-plan-meta-value`),
			),
			Div().Class(`quote-plan-meta-item`).Wrap(
				Span(`NC`).Class(`quote-plan-meta-label`),
				Span(PriceTextWholeEuro(x.noClaims, true)).Class(`quote-plan-meta-value`),
			),
			Div().Class(`quote-plan-meta-item`).Wrap(
				Span(`Comm`).Class(`quote-plan-meta-label`),
				Span(PriceTextWholeEuro(x.commission, true)).Class(`quote-plan-meta-value`),
			),
		),
		Elem(`details`).Class(`quote-plan-addon-details`).KV(`data-details-key`, Str(`plan-`, x.planId)).Wrap(
			Elem(`summary`).Class(`quote-plan-addon-title`).Wrap(
				Span(`Product details`).Class(`quote-plan-addon-title-label`),
				Span(PriceText(x.base, true)).Class(`quote-plan-addon-title-sum`),
				Span(PriceText(x.surcharge, true)).Class(`quote-plan-addon-title-sum`),
			),
			Div().Class(`quote-plan-addon-list`).Wrap(addons),
		),
	)
}

func QuoteSelectedPlanAddonView(itemId int, x QuotePlanAddon_t) Elem_t {
	pick := QuoteAddonPickText(x)
	pickView := any(pick)
	if x.hasMulti && len(x.choices) > 0 {
		key := QuoteSelectedCatKey(itemId, x.categId)
		pickView = QuotePlanAddonSelectNamedView(key, x)
	}
	return Div().Class(`quote-plan-addon`).Wrap(
		Div(x.categ).Class(`quote-plan-addon-cat`),
		Div().Class(`quote-plan-addon-pick`).Wrap(pickView),
		Div(PriceText(x.base, x.priceOk)).Class(`quote-plan-addon-base`),
		Div(PriceText(x.surcharge, x.priceOk)).Class(`quote-plan-addon-surch`),
	)
}

func QuoteSelectedPlanCardView(item QuoteSelectedItem_t, row QuotePlan_t) Elem_t {
	var addons []Elem_t
	addons = append(addons, QuotePlanBaseView(row))
	for _, addon := range row.addons { addons = append(addons, QuoteSelectedPlanAddonView(item.itemId, addon)) }

	return Div().
		Id(Str(`QuoteSelected-`, item.itemId)).
		Class(`quote-plan-card`, `quote-plan-card-selected`).
		KV(`data-selected-item-id`, Str(item.itemId)).
		KV(`data-plan-id`, Str(item.planId)).
		Wrap(
		Div().Class(`quote-plan-head`).Wrap(
			Div().Class(`quote-plan-head-main`).Wrap(
				Elem(`button`).Type(`button`).Name(QuoteSelectedDelControlName(item.itemId)).Value(Str(item.itemId)).Class(`quote-plan-pick-btn`, `quote-plan-pick-del`).Text(`🗑`),
				Div(row.label).Class(`quote-plan-label`),
			),
			Div().Class(`quote-plan-total`).Wrap(
				Div(PriceTextWholeEuro(row.price, true)).Class(`quote-plan-price`),
			),
		),
		Div().Class(`quote-plan-meta`).Wrap(
			Div().Class(`quote-plan-meta-item`).Wrap(
				Span(`Ded`).Class(`quote-plan-meta-label`),
				Span(PriceTextWholeEuro(row.deduct, true)).Class(`quote-plan-meta-value`),
			),
			Div().Class(`quote-plan-meta-item`).Wrap(
				Span(`NC`).Class(`quote-plan-meta-label`),
				Span(PriceTextWholeEuro(row.noClaims, true)).Class(`quote-plan-meta-value`),
			),
			Div().Class(`quote-plan-meta-item`).Wrap(
				Span(`Comm`).Class(`quote-plan-meta-label`),
				Span(PriceTextWholeEuro(row.commission, true)).Class(`quote-plan-meta-value`),
			),
		),
		Elem(`details`).Class(`quote-plan-addon-details`).KV(`data-details-key`, Str(`sel-`, item.itemId)).Wrap(
			Elem(`summary`).Class(`quote-plan-addon-title`).Wrap(
				Span(`Product details`).Class(`quote-plan-addon-title-label`),
				Span(PriceText(row.base, true)).Class(`quote-plan-addon-title-sum`),
				Span(PriceText(row.surcharge, true)).Class(`quote-plan-addon-title-sum`),
			),
			Div().Class(`quote-plan-addon-list`).Wrap(addons),
		),
	)
}

func QuoteFilteredPlanView(x QuotePlanFiltered_t) Elem_t {
	return Div().Class(`quote-filtered-plan`).Wrap(
		Div(x.label).Class(`quote-filtered-label`),
		Div(Join(x.reasons, `, `)).Class(`quote-filtered-reasons`),
	)
}

func QuotePlanDesktopCategs() []Categ_t {
	var out []Categ_t
	for _, categ := range App.lookup.categs.All() {
		if categ.categId == 0 { continue }
		if categ.display == 0 { continue }
		out = append(out, categ)
	}
	return out
}

func QuoteSortSelectView(sortBy string) Elem_t {
	mode := QuoteSortMode(sortBy)
	return Select(
		Option().KV(`value`, sortByName).Text(`Name`),
		Option().KV(`value`, sortByPrice).Text(`Total`),
	).Name(`sortBy`).Choose(mode).Class(`quote-plan-sort-input`)
}

func QuoteEditQuoteButton(class ...string) Elem_t {
	return Elem(`button`).
		Type(`submit`).
		KV(`formaction`, `/quote-review`).
		KV(`formmethod`, `post`).
		Class(`quote-edit-quote-btn`).
		Class(class...).
		Text(`Quote Review`)
}

func QuoteResetControlName() string { return `quoteReset` }

func QuoteResetButton(class ...string) Elem_t {
	return Elem(`button`).
		Type(`button`).
		Name(QuoteResetControlName()).
		Value(`1`).
		Class(`quote-edit-quote-btn`).
		Class(class...).
		Text(`Reset`)
}

func QuotePlanDesktopCategWidth(categ Categ_t) int {
	name := Lower(Trim(categ.name))
	switch {
	case Contains(name, `sick`):
		return 72
	case Contains(name, `pvn`):
		return 50
	case Contains(name, `hospital`):
		return 108
	case Contains(name, `dental`):
		return 108
	case Contains(name, `consumer`):
		return 84
	case Contains(name, `natural`):
		return 80
	case Contains(name, `special`):
		return 90
	}
	return 120
}

func QuotePxToRem(px int) string {
	s := fmt.Sprintf(`%.4f`, float64(px)/16.0)
	s = strings.TrimRight(s, `0`)
	s = strings.TrimRight(s, `.`)
	return s + `rem`
}

func QuotePlanDesktopGridStyle(categs []Categ_t, showVision bool) string {
	const actionColPx = 34
	const moneyColPx = 67
	x := Str(`grid-template-columns: `, QuotePxToRem(actionColPx), ` `, QuotePxToRem(moneyColPx), ` `, QuotePxToRem(moneyColPx), ` `, QuotePxToRem(moneyColPx), ` `, QuotePxToRem(294))
	for _, categ := range categs {
		x += Str(` `, QuotePxToRem(QuotePlanDesktopCategWidth(categ)))
	}
	if showVision {
		x += Str(` `, QuotePxToRem(60))
	}
	x += Str(` `, QuotePxToRem(moneyColPx))
	x += `;`
	return x
}

func QuotePlanDesktopHead(categs []Categ_t, title, sortBy string, showSort, showAction bool, action Elem_t) Elem_t {
	var cols []Elem_t
	cols = append(cols,
		Div(``).Class(`quote-plan-cell`, `quote-plan-action-head`),
		Div(`Ded`).Class(`quote-plan-cell`, `quote-plan-col-money`, `quote-plan-money-head-right`),
		Div(`NCB`).Class(`quote-plan-cell`, `quote-plan-col-money`, `quote-plan-money-head-right`),
		Div(`Total`).Class(`quote-plan-cell`, `quote-plan-col-total`, `quote-plan-money-head-right`),
	)

	var headParts []Elem_t
	headParts = append(headParts, Span(title).Class(`quote-plan-head-plan-title`))
	if showSort { headParts = append(headParts, QuoteSortSelectView(sortBy)) }
	if showAction { headParts = append(headParts, action) }
	head := Div().Class(`quote-plan-cell`, `quote-plan-name-cell`).Wrap(
		Div().Class(`quote-plan-head-plan`).Wrap(headParts),
	)
	cols = append(cols, head)

	for _, categ := range categs {
		cols = append(cols, Div(categ.name).Class(`quote-plan-cell`))
	}
	cols = append(cols, Div(`👓`).Class(`quote-plan-cell`, `quote-plan-col-money`, `quote-plan-money-head-right`))
	cols = append(cols, Div(`Comm`).Class(`quote-plan-cell`, `quote-plan-col-money`, `quote-plan-money-head-right`))
	return Div().
		Class(`quote-plan-table-row`, `quote-plan-table-head`).
		Wrap(cols)
}

func QuotePlanActionButton(name, value, icon, class string) Elem_t {
	return Elem(`button`).
		Type(`button`).
		Name(name).
		Value(value).
		Class(`quote-plan-pick-btn`, `quote-plan-desk-btn`, class).
		Text(icon)
}

func QuotePlanDesktopAddonPickView(planId int, addon QuotePlanAddon_t) Elem_t {
	return QuotePlanDesktopAddonPickNamedView(QuotePlanCatControlName(planId, addon.categId), addon)
}

func QuotePlanDesktopAddonPickNamedView(name string, addon QuotePlanAddon_t) Elem_t {
	if addon.hasMulti && len(addon.choices) > 0 {
		return Div().Class(`quote-plan-cell-pick`).Wrap(QuotePlanAddonSelectNamedView(name, addon))
	}
	return Div(QuoteAddonPickText(addon)).Class(`quote-plan-cell-pick`)
}

func QuotePlanNoAddonLevelLabel(categId CategId_t) string {
	if !IsHospDental(categId) { return `` }
	levelId := int(categId) * 10
	level, ok := App.lookup.levels.byId[levelId]
	if !ok { return `` }
	label := Trim(level.label)
	if label == `` { return `` }
	return label
}

func QuotePlanDesktopCategCellView(x QuotePlan_t, categId CategId_t) Elem_t {
	addon, ok := QuotePlanAddonByCateg(x, categId)
	if !ok {
		if label := QuotePlanNoAddonLevelLabel(categId); label != `` { return Div(label).Class(`quote-plan-cell-pick`) }
		return Div(`&nbsp;`).Class(`quote-plan-cell-pick`)
	}
	if !addon.priceOk && addon.addon == 0 && addon.level == 0 && addon.label == `` {
		if label := QuotePlanNoAddonLevelLabel(categId); label != `` { return Div(label).Class(`quote-plan-cell-pick`) }
		return Div(`&nbsp;`).Class(`quote-plan-cell-pick`)
	}
	return QuotePlanDesktopAddonPickView(x.planId, addon)
}

func QuotePlanDesktopVisionCellView(x QuotePlan_t) Elem_t {
	addon, ok := QuotePlanAddonByTag(x, `vision`)
	if !ok { return Div(`&nbsp;`).Class(`quote-plan-cell-pick`) }
	if !addon.priceOk { return Div(`&nbsp;`).Class(`quote-plan-cell-pick`) }
	amount := addon.base + addon.surcharge
	if amount == 0 { return Div(`&nbsp;`).Class(`quote-plan-cell-pick`) }
	return Div(PriceText(amount, addon.priceOk)).Class(`quote-plan-cell-pick`, `quote-plan-cell-money`, `quote-plan-cell-money-right`)
}

func QuotePlanDesktopSelectedCategCellView(itemId int, x QuotePlan_t, categId CategId_t) Elem_t {
	addon, ok := QuotePlanAddonByCateg(x, categId)
	if !ok {
		if label := QuotePlanNoAddonLevelLabel(categId); label != `` { return Div(label).Class(`quote-plan-cell-pick`) }
		return Div(`&nbsp;`).Class(`quote-plan-cell-pick`)
	}
	if !addon.priceOk && addon.addon == 0 && addon.level == 0 && addon.label == `` {
		if label := QuotePlanNoAddonLevelLabel(categId); label != `` { return Div(label).Class(`quote-plan-cell-pick`) }
		return Div(`&nbsp;`).Class(`quote-plan-cell-pick`)
	}
	name := QuoteSelectedCatKey(itemId, categId)
	return QuotePlanDesktopAddonPickNamedView(name, addon)
}

func QuotePlanDesktopSelectedAmountCategCellView(x QuotePlan_t, categId CategId_t, showBase bool) Elem_t {
	addon, ok := QuotePlanAddonByCateg(x, categId)
	if !ok { return Div(`&nbsp;`).Class(`quote-plan-cell-pick`) }
	if !addon.priceOk && addon.addon == 0 && addon.level == 0 && addon.label == `` {
		return Div(`&nbsp;`).Class(`quote-plan-cell-pick`)
	}
	amount := addon.surcharge
	if showBase { amount = addon.base }
	return Div(PriceText(amount, addon.priceOk)).Class(`quote-plan-cell-pick`, `quote-plan-cell-money`)
}

func QuotePlanDesktopSelectedAmountVisionCellView(x QuotePlan_t, showBase bool) Elem_t {
	addon, ok := QuotePlanAddonByTag(x, `vision`)
	if !ok { return Div(`&nbsp;`).Class(`quote-plan-cell-pick`) }
	if !addon.priceOk { return Div(`&nbsp;`).Class(`quote-plan-cell-pick`) }
	amount := addon.surcharge
	if showBase { amount = addon.base }
	return Div(PriceText(amount, addon.priceOk)).Class(`quote-plan-cell-pick`, `quote-plan-cell-money`)
}

func QuotePlanDesktopSelectedAmountRow(row QuotePlan_t, categs []Categ_t, showVision, showBase bool) Elem_t {
	total := row.surcharge
	planAmount := row.planSurcharge
	if showBase {
		total = row.base
		planAmount = row.planBase
	}

	var cols []Elem_t
	cols = append(cols,
		Div(`&nbsp;`).Class(`quote-plan-cell`, `quote-plan-action-cell`, `quote-plan-selected-detail-empty`),
		Div(PriceText(total, true)).Class(`quote-plan-cell`, `quote-plan-total-cell`),
		Div(`&nbsp;`).Class(`quote-plan-cell`, `quote-plan-cell-money`, `quote-plan-cell-money-right`, `quote-plan-selected-detail-empty`),
		Div(`&nbsp;`).Class(`quote-plan-cell`, `quote-plan-cell-money`, `quote-plan-cell-money-right`, `quote-plan-selected-detail-empty`),
		Div().Class(`quote-plan-cell`, `quote-plan-name-cell`, `quote-plan-selected-detail-name`).Wrap(
			Div().Class(`quote-plan-selected-detail-split`).Wrap(
				Div(PriceText(planAmount, row.planOk)).Class(`quote-plan-selected-detail-split-left`),
				Div(`&nbsp;`).Class(`quote-plan-selected-detail-split-right`, `quote-plan-selected-detail-empty`),
			),
		),
	)
	for _, categ := range categs {
		cols = append(cols, Div().Class(`quote-plan-cell`).Wrap(QuotePlanDesktopSelectedAmountCategCellView(row, categ.categId, showBase)))
	}
	if showVision {
		cols = append(cols, Div().Class(`quote-plan-cell`).Wrap(QuotePlanDesktopSelectedAmountVisionCellView(row, showBase)))
	}
	cols = append(cols, Div(`&nbsp;`).Class(`quote-plan-cell`, `quote-plan-cell-money`, `quote-plan-cell-money-right`, `quote-plan-selected-detail-empty`))

	rowClass := `quote-plan-table-selected-surch-row`
	if showBase { rowClass = `quote-plan-table-selected-base-row` }
	return Div().
		Class(`quote-plan-table-row`, `quote-plan-table-selected-detail-row`, rowClass).
		KV(`style`, QuotePlanDesktopGridStyle(categs, showVision)).
		Wrap(cols)
}

func PriceTextWholeEuro(amount EuroCent_t, ok bool) string {
	if !ok { return `-` }
	return EuroFlatFromCent(amount).OutEuro()
}

func QuotePlanDesktopRow(x QuotePlan_t, categs []Categ_t) Elem_t {
	var cols []Elem_t
	cols = append(cols,
		Div().Class(`quote-plan-cell`, `quote-plan-action-cell`).Wrap(
			QuotePlanActionButton(QuoteSelectedAddControlName(x.planId), Str(x.planId), `🛒`, `quote-plan-pick-add`),
		),
		Div(PriceTextWholeEuro(x.deduct, true)).Class(`quote-plan-cell`, `quote-plan-col-money`, `quote-plan-cell-money`, `quote-plan-cell-money-right`),
		Div(PriceTextWholeEuro(x.noClaims, true)).Class(`quote-plan-cell`, `quote-plan-col-money`, `quote-plan-cell-money`, `quote-plan-cell-money-right`),
		Div(PriceTextWholeEuro(x.price, true)).Class(`quote-plan-cell`, `quote-plan-col-total`, `quote-plan-cell-money`, `quote-plan-cell-money-right`),
		Div(x.label).Class(`quote-plan-cell`, `quote-plan-name-cell`),
	)
	for _, categ := range categs {
		cols = append(cols, Div().Class(`quote-plan-cell`).Wrap(QuotePlanDesktopCategCellView(x, categ.categId)))
	}
	cols = append(cols, Div().Class(`quote-plan-cell`, `quote-plan-col-money`).Wrap(QuotePlanDesktopVisionCellView(x)))
	cols = append(cols, Div(PriceTextWholeEuro(x.commission, true)).Class(`quote-plan-cell`, `quote-plan-col-money`, `quote-plan-cell-money`, `quote-plan-cell-money-right`))
	return Div().
		Class(`quote-plan-table-row`).
		Wrap(cols)
}

func QuotePlanDesktopSelectedRow(item QuoteSelectedItem_t, row QuotePlan_t, categs []Categ_t) Elem_t {
	var cols []Elem_t
	cols = append(cols,
		Div().Class(`quote-plan-cell`, `quote-plan-action-cell`).Wrap(
			QuotePlanActionButton(QuoteSelectedDelControlName(item.itemId), Str(item.itemId), `🗑`, `quote-plan-pick-del`),
		),
		Div(PriceTextWholeEuro(row.deduct, true)).Class(`quote-plan-cell`, `quote-plan-col-money`, `quote-plan-cell-money`, `quote-plan-cell-money-right`),
		Div(PriceTextWholeEuro(row.noClaims, true)).Class(`quote-plan-cell`, `quote-plan-col-money`, `quote-plan-cell-money`, `quote-plan-cell-money-right`),
		Div(PriceTextWholeEuro(row.price, true)).Class(`quote-plan-cell`, `quote-plan-col-total`, `quote-plan-cell-money`, `quote-plan-cell-money-right`),
		Div(row.label).Class(`quote-plan-cell`, `quote-plan-name-cell`),
	)
	for _, categ := range categs {
		cols = append(cols, Div().Class(`quote-plan-cell`).Wrap(QuotePlanDesktopSelectedCategCellView(item.itemId, row, categ.categId)))
	}
	cols = append(cols, Div().Class(`quote-plan-cell`, `quote-plan-col-money`).Wrap(QuotePlanDesktopVisionCellView(row)))
	cols = append(cols, Div(PriceTextWholeEuro(row.commission, true)).Class(`quote-plan-cell`, `quote-plan-col-money`, `quote-plan-cell-money`, `quote-plan-cell-money-right`))
	return Div().
		Class(`quote-plan-table-row`, `quote-plan-table-selected-main-row`).
		Wrap(cols)
}

func QuotePlanDesktopView(data QuotePlans_t) Elem_t {
	categs := QuotePlanDesktopCategs()
	var rows []Elem_t
	rows = append(rows, QuotePlanDesktopHead(categs, `Plan`, data.sortBy, true, false, Div()))
	for _, x := range data.plans { rows = append(rows, QuotePlanDesktopRow(x, categs)) }
	return Div().Class(`quote-plan-table`, `quote-plan-table-main`).Wrap(rows)
}

func QuoteDesktopSelectedPlansBox(vars QuoteVars_t) Elem_t {
	state := QuoteStateFromQuoteVars(vars)
	selectedRowsData := QuoteSelectedRows(state)
	categs := QuotePlanDesktopCategs()

	var selectedRows []Elem_t
	for _, x := range selectedRowsData {
		selectedRows = append(selectedRows, QuotePlanDesktopSelectedRow(x.item, x.row, categs))
	}

	showEdit := len(selectedRows) > 0
	var rows []Elem_t
	rows = append(rows, QuotePlanDesktopHead(categs, `Plan`, ``, false, showEdit, QuoteEditQuoteButton(`quote-desk-selected-edit-btn`)))
	rows = append(rows, selectedRows...)

	out := Div().Id(`QuoteDeskSelected`).Class(`quote-desk-selected`).Wrap(
		Div().Class(`quote-plan-table`, `quote-plan-table-selected`).Wrap(rows),
	)
	if len(selectedRows) == 0 {
		out = Div().Id(`QuoteDeskSelected`).Class(`quote-desk-selected`).Wrap(
			Div().Class(`quote-plan-table`, `quote-plan-table-selected`).Wrap(rows),
			Div(`No plans selected.`).Class(`quote-desk-selected-empty`),
		)
	}
	return out
}

func QuotePhoneSelectedPlansBox(vars QuoteVars_t) Elem_t {
	state := QuoteStateFromQuoteVars(vars)
	selectedRowsData := QuoteSelectedRows(state)
	var cards []Elem_t
	for _, x := range selectedRowsData {
		cards = append(cards, QuoteSelectedPlanCardView(x.item, x.row))
	}
	if len(cards) == 0 {
		cards = append(cards, Div(`No plans selected.`).Class(`quote-selected-empty`))
	}
	var titleBar []Elem_t
	titleBar = append(titleBar, Span(QuoteSelectedTitle(len(selectedRowsData))))
	if len(selectedRowsData) > 0 {
		titleBar = append(titleBar, QuoteEditQuoteButton(`quote-phone-selected-edit-btn`))
	}
	return Elem(`details`).Id(`QuoteSelectedCard`).Class(`quote-card`, `quote-phone-card`, `quote-phone-fold`, `quote-phone-selected-fold`, `quote-phone-selected-card`).Wrap(
		Elem(`summary`).Class(`quote-card-title`, `quote-phone-fold-title`, `quote-phone-selected-title`).Wrap(titleBar),
		Div().Class(`quote-phone-selected-list`).Wrap(cards),
	)
}

func QuoteFilteredPlansBox(filteredPlans []QuotePlanFiltered_t) Elem_t {
	var filtered []Elem_t
	for _, x := range filteredPlans { filtered = append(filtered, QuoteFilteredPlanView(x)) }
	if len(filtered) == 0 {
		filtered = append(filtered, Div(`No plans filtered out.`).Class(`quote-filtered-empty`))
	}
	return Elem(`details`).Class(`quote-filter-box`).Wrap(
		Elem(`summary`).Class(`quote-filter-box-title`).Wrap(
			Span(`Filtered plans (` , len(filteredPlans), `)`),
		),
		Div().Class(`quote-filtered-list`).Wrap(filtered),
	)
}
