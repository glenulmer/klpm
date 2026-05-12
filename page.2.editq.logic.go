package main

import (
	"fmt"
	"sort"
	"strings"

	. "klpm/lib/date"
	. "klpm/lib/dec2"
	. "klpm/lib/output"
)

const editQDepMaxCount = 5
const editQPreexModePct = `pct`
const editQPreexModeEur = `eur`

type EditQDep_t struct {
	depId int
	name string
	birth string
	vision bool
}

type EditQPreexCharge_t struct {
	itemId int
	plan string
	planPrice EuroCent_t
	planAge int
	planAgeMode string
	categId CategId_t
	level string
	base EuroCent_t
	mode string
	amount string
	note string
	applied EuroCent_t
}

func EditQDepNameKey(depId int) string { return Str(`editq-dep-`, depId, `-name`) }
func EditQDepBirthKey(depId int) string { return Str(`editq-dep-`, depId, `-birth`) }
func EditQDepVisionKey(depId int) string { return Str(`editq-dep-`, depId, `-vision`) }

func EditQDepNameControl(name string) (depId int, ok bool) {
	n, err := fmt.Sscanf(name, `editq-dep-%d-name`, &depId)
	if err != nil || n != 1 || depId <= 0 { return 0, false }
	return depId, true
}

func EditQDepBirthControl(name string) (depId int, ok bool) {
	n, err := fmt.Sscanf(name, `editq-dep-%d-birth`, &depId)
	if err != nil || n != 1 || depId <= 0 { return 0, false }
	return depId, true
}

func EditQDepVisionControl(name string) (depId int, ok bool) {
	n, err := fmt.Sscanf(name, `editq-dep-%d-vision`, &depId)
	if err != nil || n != 1 || depId <= 0 { return 0, false }
	return depId, true
}

func EditQDepAddControlName() string { return `editq-dep-add` }
func EditQDepAddControl(name string) bool { return name == EditQDepAddControlName() }

func EditQDepDelControlName(depId int) string { return Str(`editq-dep-del-`, depId) }
func EditQDepDelControl(name string) (depId int, ok bool) {
	n, err := fmt.Sscanf(name, `editq-dep-del-%d`, &depId)
	if err != nil || n != 1 || depId <= 0 { return 0, false }
	return depId, true
}

func EditQPreexModeKey(itemId int, categId CategId_t) string {
	return Str(`editq-preex-`, itemId, `-`, categId, `-mode`)
}

func EditQPreexAmountKey(itemId int, categId CategId_t) string {
	return Str(`editq-preex-`, itemId, `-`, categId, `-amount`)
}

func EditQPreexNoteKey(itemId int, categId CategId_t) string {
	return Str(`editq-preex-`, itemId, `-`, categId, `-note`)
}

func EditQPreexModeControl(name string) (itemId int, categId CategId_t, ok bool) {
	var cat int
	n, err := fmt.Sscanf(name, `editq-preex-%d-%d-mode`, &itemId, &cat)
	if err != nil || n != 2 || itemId <= 0 || cat < 0 { return 0, 0, false }
	return itemId, CategId_t(cat), true
}

func EditQPreexAmountControl(name string) (itemId int, categId CategId_t, ok bool) {
	var cat int
	n, err := fmt.Sscanf(name, `editq-preex-%d-%d-amount`, &itemId, &cat)
	if err != nil || n != 2 || itemId <= 0 || cat < 0 { return 0, 0, false }
	return itemId, CategId_t(cat), true
}

func EditQPreexNoteControl(name string) (itemId int, categId CategId_t, ok bool) {
	var cat int
	n, err := fmt.Sscanf(name, `editq-preex-%d-%d-note`, &itemId, &cat)
	if err != nil || n != 2 || itemId <= 0 || cat < 0 { return 0, 0, false }
	return itemId, CategId_t(cat), true
}

func EditQDepChargeModeKey(depId, itemId int, categId CategId_t) string {
	return Str(`editq-dep-`, depId, `-preex-`, itemId, `-`, categId, `-mode`)
}

func EditQDepChargeAmountKey(depId, itemId int, categId CategId_t) string {
	return Str(`editq-dep-`, depId, `-preex-`, itemId, `-`, categId, `-amount`)
}

func EditQDepChargeNoteKey(depId, itemId int, categId CategId_t) string {
	return Str(`editq-dep-`, depId, `-preex-`, itemId, `-`, categId, `-note`)
}

func EditQDepChargeModeControl(name string) (depId, itemId int, categId CategId_t, ok bool) {
	var cat int
	n, err := fmt.Sscanf(name, `editq-dep-%d-preex-%d-%d-mode`, &depId, &itemId, &cat)
	if err != nil || n != 3 || depId <= 0 || itemId <= 0 || cat < 0 { return 0, 0, 0, false }
	return depId, itemId, CategId_t(cat), true
}

func EditQDepChargeAmountControl(name string) (depId, itemId int, categId CategId_t, ok bool) {
	var cat int
	n, err := fmt.Sscanf(name, `editq-dep-%d-preex-%d-%d-amount`, &depId, &itemId, &cat)
	if err != nil || n != 3 || depId <= 0 || itemId <= 0 || cat < 0 { return 0, 0, 0, false }
	return depId, itemId, CategId_t(cat), true
}

func EditQDepChargeNoteControl(name string) (depId, itemId int, categId CategId_t, ok bool) {
	var cat int
	n, err := fmt.Sscanf(name, `editq-dep-%d-preex-%d-%d-note`, &depId, &itemId, &cat)
	if err != nil || n != 3 || depId <= 0 || itemId <= 0 || cat < 0 { return 0, 0, 0, false }
	return depId, itemId, CategId_t(cat), true
}

func EditQPreexMode(v string) string {
	if Lower(Trim(v)) == editQPreexModeEur { return editQPreexModeEur }
	return editQPreexModePct
}

func EditQParseDecimal100(v string) (int64, bool) {
	s := Trim(v)
	s = strings.ReplaceAll(s, `%`, ``)
	s = strings.ReplaceAll(s, `€`, ``)
	s = strings.ReplaceAll(s, `eur`, ``)
	s = strings.ReplaceAll(s, ` `, ``)
	s = strings.ReplaceAll(s, `.`, `,`)
	if s == `` { return 0, false }
	parts := Split(s, `,`)
	if len(parts) > 2 { return 0, false }
	onlyDigits := func(in string) string {
		var out []byte
		for i := 0; i < len(in); i++ {
			if in[i] >= '0' && in[i] <= '9' { out = append(out, in[i]) }
		}
		return string(out)
	}
	whole := onlyDigits(parts[0])
	if whole == `` { whole = `0` }
	frac := ``
	if len(parts) == 2 { frac = onlyDigits(parts[1]) }
	if len(frac) > 2 { frac = frac[:2] }
	for len(frac) < 2 { frac += `0` }
	return int64(Atoi(whole)*100 + Atoi(frac)), true
}

func EditQPreexAppliedAmount(mode, amount string, base EuroCent_t) EuroCent_t {
	n, ok := EditQParseDecimal100(amount)
	if !ok || n <= 0 { return 0 }
	if EditQPreexMode(mode) == editQPreexModeEur {
		return EuroCent_t(n)
	}
	// n is percent*100, apply to base euro-cent amount.
	return EuroCent_t((int64(base) * n) / 10000)
}

func EditQChoicePreex(choice PlanQuoteInfo_t, categId CategId_t) (Preex_t, bool) {
	for _, preex := range choice.preex {
		if preex.categ != categId { continue }
		return preex, true
	}
	return Preex_t{}, false
}

func EditQPreexFields(preex Preex_t) (mode, amount, note string) {
	mode, amount = QuotePreexModeAmount(preex)
	note = preex.note
	return mode, amount, note
}

func EditQPreexCharges(vars QuoteVars_t) []EditQPreexCharge_t {
	state := QuoteStateFromQuoteVars(vars)
	selected := QuoteSelectedRows(state)
	var out []EditQPreexCharge_t
	for _, x := range selected {
		if x.row.planOk && x.row.planBase+x.row.planSurcharge > 0 {
			mode, amount, note := editQPreexModePct, ``, ``
			if choice, ok := vars.choices[ChoiceId_t(x.item.itemId)]; ok {
				if preex, has := EditQChoicePreex(choice, 0); has {
					mode, amount, note = EditQPreexFields(preex)
				}
			}
			applied := EditQPreexAppliedAmount(mode, amount, x.row.base)
			out = append(out, EditQPreexCharge_t{
				itemId: x.item.itemId,
				plan: x.row.label,
				planPrice: x.row.price,
				categId: 0,
				level: `Plan`,
				base: x.row.base,
				mode: mode,
				amount: amount,
				note: note,
				applied: applied,
			})
		}

		for _, addon := range x.row.addons {
			if !addon.priceOk { continue }
			if addon.categId <= 0 { continue }
			if Contains(Lower(Trim(addon.categ)), `vision`) { continue }
			if addon.base+addon.surcharge <= 0 { continue }
			mode, amount, note := editQPreexModePct, ``, ``
			if choice, ok := vars.choices[ChoiceId_t(x.item.itemId)]; ok {
				if preex, has := EditQChoicePreex(choice, addon.categId); has {
					mode, amount, note = EditQPreexFields(preex)
				}
			}
			applied := EditQPreexAppliedAmount(mode, amount, addon.base)
			level := QuoteAddonPickText(addon)
			if level == `` { level = addon.categ }
			out = append(out, EditQPreexCharge_t{
				itemId: x.item.itemId,
				plan: x.row.label,
				planPrice: x.row.price,
				categId: addon.categId,
				level: level,
				base: addon.base,
				mode: mode,
				amount: amount,
				note: note,
				applied: applied,
			})
		}
	}
	return out
}

func EditQDepChoicePreex(dep Dependant_t, choiceId ChoiceId_t, categId CategId_t) (Preex_t, bool) {
	list := dep.preexByChoice[choiceId]
	for _, preex := range list {
		if preex.categ != categId { continue }
		return preex, true
	}
	return Preex_t{}, false
}

func EditQDependantCharges(vars QuoteVars_t, dep EditQDep_t) []EditQPreexCharge_t {
	work := QuoteStateFromQuoteVars(vars)
	work.quote.core.birth = QuoteParseBirthDate(dep.birth)
	work.quote.core.vision = dep.vision
	_, yearAge, exactAge := PlanAges(work)
	selected := QuoteSelectedRows(work)
	depData := Dependant_t{}
	for _, x := range vars.dependants {
		if x.depId != dep.depId { continue }
		depData = x
		break
	}

	var out []EditQPreexCharge_t
	for _, x := range selected {
		age := yearAge
		ageMode := `year`
		if EditQPlanUsesExactAge(x.row.planId) {
			age = exactAge
			ageMode = `exact`
		}
		if x.row.planOk && x.row.planBase+x.row.planSurcharge > 0 {
			mode, amount, note := editQPreexModePct, ``, ``
			if preex, has := EditQDepChoicePreex(depData, ChoiceId_t(x.item.itemId), 0); has {
				mode, amount, note = EditQPreexFields(preex)
			}
			applied := EditQPreexAppliedAmount(mode, amount, x.row.base)
			out = append(out, EditQPreexCharge_t{
				itemId: x.item.itemId,
				plan: x.row.label,
				planPrice: x.row.price,
				planAge: age,
				planAgeMode: ageMode,
				categId: 0,
				level: `Plan`,
				base: x.row.base,
				mode: mode,
				amount: amount,
				note: note,
				applied: applied,
			})
		}
		for _, addon := range x.row.addons {
			if !addon.priceOk { continue }
			if addon.categId <= 0 { continue }
			if Contains(Lower(Trim(addon.categ)), `vision`) { continue }
			if addon.base+addon.surcharge <= 0 { continue }
			mode, amount, note := editQPreexModePct, ``, ``
			if preex, has := EditQDepChoicePreex(depData, ChoiceId_t(x.item.itemId), addon.categId); has {
				mode, amount, note = EditQPreexFields(preex)
			}
			applied := EditQPreexAppliedAmount(mode, amount, addon.base)
			level := QuoteAddonPickText(addon)
			if level == `` { level = addon.categ }
			out = append(out, EditQPreexCharge_t{
				itemId: x.item.itemId,
				plan: x.row.label,
				planPrice: x.row.price,
				planAge: age,
				planAgeMode: ageMode,
				categId: addon.categId,
				level: level,
				base: addon.base,
				mode: mode,
				amount: amount,
				note: note,
				applied: applied,
			})
		}
	}
	return out
}

func EditQPlanUsesExactAge(planId int) bool {
	for id, plan := range App.lookup.plans.All() {
		if id != planId { continue }
		return plan.exactAge
	}
	return false
}

func EditQDefaultDependantBirth(vars QuoteVars_t) string {
	year := 0
	buy := vars.core.buy
	if buy.Year() > 0 { year = buy.Year() }
	if year <= 0 {
		_, _, def := QuoteBuyBounds()
		if def.Year() > 0 { year = def.Year() }
	}
	if year <= 0 {
		today := CurrentDBDate()
		if today.Year() > 0 { year = today.Year() }
	}
	if year <= 0 { year = 2026 }
	return fmt.Sprintf(`%04d-06-15`, year-4)
}

func EditQBirthSortKey(v string) string {
	v = Trim(v)
	if len(v) != len(`2006-01-02`) { return `9999-99-99` }
	if v[4] != '-' || v[7] != '-' { return `9999-99-99` }
	return v
}

func EditQDependants(vars QuoteVars_t, sortForGet bool) []EditQDep_t {
	var out []EditQDep_t
	for i, dep := range vars.dependants {
		depId := dep.depId
		if depId <= 0 { depId = i + 1 }
		birth := ``
		if Valid(dep.birth) { birth = dep.birth.Format(`yyyy-mm-dd`) }
		out = append(out, EditQDep_t{
			depId: depId,
			name: dep.name,
			birth: birth,
			vision: dep.vision,
		})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].depId != out[j].depId { return out[i].depId < out[j].depId }
		return false
	})
	if !sortForGet { return out }

	sort.Slice(out, func(i, j int) bool {
		bi, bj := EditQBirthSortKey(out[i].birth), EditQBirthSortKey(out[j].birth)
		if bi != bj { return bi < bj }
		ni, nj := Lower(Trim(out[i].name)), Lower(Trim(out[j].name))
		if ni != nj { return ni < nj }
		return out[i].depId < out[j].depId
	})
	return out
}

func EditQDeleteDependant(vars *QuoteVars_t, depId int) {
	if vars == nil || depId <= 0 { return }
	var out []Dependant_t
	for _, dep := range vars.dependants {
		if dep.depId == depId { continue }
		out = append(out, dep)
	}
	vars.dependants = out
}

func EditQSetPreex(preex []Preex_t, categId CategId_t, mode, amount, note string, setMode, setAmount, setNote bool) []Preex_t {
	pick := -1
	for i, x := range preex {
		if x.categ != categId { continue }
		pick = i
		break
	}
	curr := Preex_t{ categ:categId }
	if pick >= 0 { curr = preex[pick] }

	if setMode {
		curr.mode = EditQPreexMode(mode)
		if curr.mode == editQPreexModeEur {
			if curr.amount.percent > 0 {
				curr.amount.euro = EuroCent_t(curr.amount.percent)
				curr.amount.percent = 0
			}
		} else {
			if curr.amount.euro > 0 {
				curr.amount.percent = Percent_t(curr.amount.euro)
				curr.amount.euro = 0
			}
		}
	}
	if setAmount {
		modeNow := EditQPreexMode(mode)
		if mode == `` && curr.mode != `` { modeNow = EditQPreexMode(curr.mode) }
		curr.mode = modeNow
		n, ok := EditQParseDecimal100(amount)
		if !ok || n <= 0 {
			curr.amount.percent = 0
			curr.amount.euro = 0
		} else if modeNow == editQPreexModeEur {
			curr.amount.euro = EuroCent_t(n)
			curr.amount.percent = 0
		} else {
			curr.amount.percent = Percent_t(n)
			curr.amount.euro = 0
		}
	}
	if setNote {
		curr.note = note
	}

	if pick >= 0 {
		preex[pick] = curr
		return preex
	}
	return append(preex, curr)
}

func EditQUpsertDep(vars *QuoteVars_t, depId int) int {
	if vars == nil || depId <= 0 { return -1 }
	for i, dep := range vars.dependants {
		if dep.depId != depId { continue }
		if dep.preexByChoice == nil { dep.preexByChoice = make(map[ChoiceId_t][]Preex_t) }
		vars.dependants[i] = dep
		return i
	}
	dep := Dependant_t{ depId:depId, preexByChoice:make(map[ChoiceId_t][]Preex_t) }
	vars.dependants = append(vars.dependants, dep)
	return len(vars.dependants)-1
}

func EditQSetDepName(vars *QuoteVars_t, depId int, value string) {
	if idx := EditQUpsertDep(vars, depId); idx >= 0 {
		vars.dependants[idx].name = value
	}
}

func EditQSetDepBirth(vars *QuoteVars_t, depId int, value string) {
	if idx := EditQUpsertDep(vars, depId); idx >= 0 {
		vars.dependants[idx].birth = QuoteParseBirthDate(value)
	}
}

func EditQSetDepVision(vars *QuoteVars_t, depId int, value string) {
	if idx := EditQUpsertDep(vars, depId); idx >= 0 {
		vars.dependants[idx].vision = QuoteVarBool(value)
	}
}

func EditQSetDepPreex(vars *QuoteVars_t, depId int, itemId ChoiceId_t, categId CategId_t, mode, amount, note string, setMode, setAmount, setNote bool) {
	idx := EditQUpsertDep(vars, depId)
	if idx < 0 { return }
	dep := vars.dependants[idx]
	list := dep.preexByChoice[itemId]
	list = EditQSetPreex(list, categId, mode, amount, note, setMode, setAmount, setNote)
	dep.preexByChoice[itemId] = list
	vars.dependants[idx] = dep
}

func EditQAddDependant(state *State_t) bool {
	QuoteEnsureDefaults(state)
	deps := EditQDependants(state.quote, false)
	if len(deps) >= editQDepMaxCount { return false }
	childNo := len(deps) + 1
	next := 1
	for _, dep := range deps {
		if dep.depId >= next { next = dep.depId + 1 }
	}
	state.quote.dependants = append(state.quote.dependants, Dependant_t{
		depId: next,
		name: Str(`Child `, childNo),
		birth: QuoteParseBirthDate(EditQDefaultDependantBirth(state.quote)),
		vision: false,
		preexByChoice: make(map[ChoiceId_t][]Preex_t),
	})
	return true
}

func EditQApply(state *State_t, name, value string) bool {
	if name == `` { return false }
	QuoteEnsureDefaults(state)

	if name == `clientName` {
		state.quote.core.clientName = value
		return true
	}
	if name == `lang` {
		if Atoi(value) <= 0 { value = Str(int(English)) }
		state.quote.lang = LangId_t(Atoi(value))
		return true
	}
	if name == `slim` {
		state.quote.slim = QuoteVarBool(value)
		return true
	}

	if itemId, categId, ok := EditQPreexModeControl(name); ok {
		choice := state.quote.choices[ChoiceId_t(itemId)]
		choice.preex = EditQSetPreex(choice.preex, categId, EditQPreexMode(value), ``, ``, true, false, false)
		state.quote.choices[ChoiceId_t(itemId)] = choice
		return true
	}
	if itemId, categId, ok := EditQPreexAmountControl(name); ok {
		choice := state.quote.choices[ChoiceId_t(itemId)]
		mode, _, _ := EditQPreexFields(Preex_t{})
		if current, has := EditQChoicePreex(choice, categId); has {
			mode, _, _ = EditQPreexFields(current)
		}
		choice.preex = EditQSetPreex(choice.preex, categId, mode, value, ``, false, true, false)
		state.quote.choices[ChoiceId_t(itemId)] = choice
		return true
	}
	if itemId, categId, ok := EditQPreexNoteControl(name); ok {
		choice := state.quote.choices[ChoiceId_t(itemId)]
		choice.preex = EditQSetPreex(choice.preex, categId, ``, ``, value, false, false, true)
		state.quote.choices[ChoiceId_t(itemId)] = choice
		return true
	}
	if depId, itemId, categId, ok := EditQDepChargeModeControl(name); ok {
		EditQSetDepPreex(&state.quote, depId, ChoiceId_t(itemId), categId, EditQPreexMode(value), ``, ``, true, false, false)
		return true
	}
	if depId, itemId, categId, ok := EditQDepChargeAmountControl(name); ok {
		mode := editQPreexModePct
		for _, dep := range state.quote.dependants {
			if dep.depId != depId { continue }
			if preex, has := EditQDepChoicePreex(dep, ChoiceId_t(itemId), categId); has {
				mode, _, _ = EditQPreexFields(preex)
			}
			break
		}
		EditQSetDepPreex(&state.quote, depId, ChoiceId_t(itemId), categId, mode, value, ``, false, true, false)
		return true
	}
	if depId, itemId, categId, ok := EditQDepChargeNoteControl(name); ok {
		EditQSetDepPreex(&state.quote, depId, ChoiceId_t(itemId), categId, ``, ``, value, false, false, true)
		return true
	}

	if EditQDepAddControl(name) {
		EditQAddDependant(state)
		return true
	}
	if depId, ok := EditQDepDelControl(name); ok {
		EditQDeleteDependant(&state.quote, depId)
		return true
	}

	if depId, ok := EditQDepNameControl(name); ok {
		EditQSetDepName(&state.quote, depId, value)
		return true
	}
	if depId, ok := EditQDepBirthControl(name); ok {
		EditQSetDepBirth(&state.quote, depId, value)
		return true
	}
	if depId, ok := EditQDepVisionControl(name); ok {
		EditQSetDepVision(&state.quote, depId, value)
		return true
	}
	return false
}

func EditQDropPreinsertedDependant(state *State_t) {
	if state == nil { return }
	QuoteEnsureDefaults(state)
	deps := EditQDependants(state.quote, false)
	if len(deps) != 1 { return }

	dep := deps[0]
	if dep.depId <= 0 { return }
	if Trim(dep.name) != `` || dep.vision { return }
	if Trim(dep.birth) != EditQDefaultDependantBirth(state.quote) { return }
	for _, x := range state.quote.dependants {
		if x.depId != dep.depId { continue }
		for _, list := range x.preexByChoice {
			for _, preex := range list {
				if preex.note != `` { return }
				if preex.amount.euro > 0 || preex.amount.percent > 0 { return }
			}
		}
		break
	}
	EditQDeleteDependant(&state.quote, dep.depId)
}
