package main

import (
	"fmt"
	"sort"

	. "klpm/lib/date"
	. "klpm/lib/dec2"
	. "klpm/lib/output"
)

type QuotePlanAddon_t struct {
	categId CategId_t
	categ string
	label string
	addon AddonId_t
	level int
	hasMulti bool
	base EuroCent_t
	surcharge EuroCent_t
	priceOk bool
	choices []CatChoice_t
}

type QuotePlan_t struct {
	planId int
	label string
	price EuroCent_t
	deduct EuroCent_t
	noClaims EuroCent_t
	commission EuroCent_t
	planBase EuroCent_t
	planSurcharge EuroCent_t
	planOk bool
	base EuroCent_t
	surcharge EuroCent_t
	addons []QuotePlanAddon_t
}

type QuotePlanFiltered_t struct {
	planId int
	label string
	reasons []string
}

type QuotePlans_t struct {
	plans []QuotePlan_t
	filtered []QuotePlanFiltered_t
	showVision bool
	sortBy string
}

type PlanFilters_t struct {
	segment int
	prior int
	noExam bool
	specref int
	vision bool
	tempVisa bool
	noPVN bool
	naturalMed bool
	deductMin, deductMax int
	hospitalMin, hospitalMax int
	dentalMin, dentalMax int
}

func QuotePlanCatControlName(planId int, categId CategId_t) string {
	return Str(`plancat-`, planId, `-`, categId)
}

func QuotePlanCatControl(name string) (planId int, categId CategId_t, ok bool) {
	var cat int
	n, err := fmt.Sscanf(name, `plancat-%d-%d`, &planId, &cat)
	if err != nil || n != 2 { return 0, 0, false }
	if planId <= 0 || cat <= 0 { return 0, 0, false }
	return planId, CategId_t(cat), true
}

func QuoteChoiceByAddon(choices []CatChoice_t, addon AddonId_t) (CatChoice_t, bool) {
	for _, choice := range choices {
		if choice.addon == addon { return choice, true }
	}
	return CatChoice_t{}, false
}

func PlanAges(state State_t) (buyYear, yearAge, exactAge int) {
	buy := StateDate(state, `buy`)
	birth := StateDate(state, `birth`)
	if !Valid(buy) || !Valid(birth) { return 0, 0, 0 }

	buyYear = buy.Year()
	yearAge = buyYear - birth.Year()
	exactAge = yearAge
	if int(buy)%10000 < int(birth)%10000 { exactAge-- }
	return buyYear, yearAge, exactAge
}

func LookupPrice(buyYear, age, productId int) (base, surcharge EuroCent_t, ok bool) {
	price, ok := App.lookup.prices[YAP(buyYear, age, productId)]
	if !ok { return 0, 0, false }
	return price.base, price.surcharge, true
}

func CatPrice(buyYear, age, productId int, categId CategId_t, sickCover int, suppressSurch bool) (base, surcharge EuroCent_t, ok bool) {
	base, surcharge, ok = LookupPrice(buyYear, age, productId)
	if !ok { return 0, 0, false }
	if categId == catSick {
		base *= EuroCent_t(sickCover / 4500)
	}
	if suppressSurch { surcharge = 0 }
	return base, surcharge, true
}

func PriceText(amount EuroCent_t, ok bool) string {
	if !ok { return `-` }
	return amount.OutEuro()
}

func AddonName(choice CatChoice_t) string {
	label := Trim(choice.label)
	if label != `` { return label }
	product, ok := App.lookup.products[ProductId_t(choice.addon)]
	if ok { return product.name }
	return Str(choice.addon)
}

func CategIs(categ Categ_t, tag string) bool {
	return Contains(Lower(Trim(categ.name)), Lower(Trim(tag)))
}

func PickAnyAddon(choices []CatChoice_t) AddonId_t {
	for _, choice := range choices {
		if choice.addon != 0 { return choice.addon }
	}
	return 0
}

func IsHospDental(categId CategId_t) bool {
	return int(categId) == catHospital || int(categId) == catDental
}

func QuoteVars(state *State_t) QuoteVars_t {
	if state == nil { return QuoteDefaultVars() }
	QuoteEnsureDefaults(state)
	return CloneQuoteVars(state.quote)
}

func QuotePreexModeAmount(preex Preex_t) (mode, amount string) {
	if EditQPreexMode(preex.mode) == editQPreexModeEur {
		if preex.amount.euro > 0 { return editQPreexModeEur, preex.amount.euro.String() }
		return editQPreexModeEur, ``
	}
	if preex.amount.euro > 0 { return editQPreexModeEur, preex.amount.euro.String() }
	if preex.amount.percent > 0 { return editQPreexModePct, EuroCent_t(preex.amount.percent).String() }
	return editQPreexModePct, ``
}

func QuoteStateFromQuoteVars(vars QuoteVars_t) State_t {
	out := InitState()
	out.quote = CloneQuoteVars(vars)
	QuoteEnsureVars(&out.quote)
	return out
}

func PlanBundledLevel(plan Plan_t, categId CategId_t) int {
	switch int(categId) {
	case catHospital:
		return plan.hospital
	case catDental:
		return plan.dental
	}
	return 0
}

func LevelLabel(level int) string {
	x, ok := App.lookup.levels.byId[level]
	if !ok { return Str(level) }
	label := Trim(x.label)
	if label != `` { return label }
	return Str(level)
}

func EnsureFreeChoice(plan Plan_t, categId CategId_t, choices []CatChoice_t) []CatChoice_t {
	if !IsHospDental(categId) { return choices }
	for _, choice := range choices {
		if choice.addon == 0 { return choices }
	}

	level := PlanBundledLevel(plan, categId)
	free := CatChoice_t{
		addon: 0,
		level: level,
		label: LevelLabel(level),
	}
	return append([]CatChoice_t{ free }, choices...)
}

func HasZeroChoice(choices []CatChoice_t) bool {
	for _, choice := range choices {
		if choice.addon == 0 { return true }
	}
	return false
}

func HasAddonChoice(choices []CatChoice_t) bool {
	for _, choice := range choices {
		if choice.addon != 0 { return true }
	}
	return false
}

func EnsureOptionalChoice(categ Categ_t, choices []CatChoice_t) []CatChoice_t {
	if categ.required != 0 { return choices }
	if HasZeroChoice(choices) || !HasAddonChoice(choices) { return choices }
	free := CatChoice_t{ addon: 0, level: 0, label: `---` }
	return append([]CatChoice_t{ free }, choices...)
}

func PickRangeChoice(plan Plan_t, categId CategId_t, choices []CatChoice_t, minLevel, maxLevel int) CatChoice_t {
	bundled := PlanBundledLevel(plan, categId)
	if bundled >= minLevel && bundled <= maxLevel {
		return CatChoice_t{ addon: 0, level: bundled, label: LevelLabel(bundled) }
	}
	for _, choice := range choices {
		if choice.addon == 0 { continue }
		if choice.level >= minLevel && choice.level <= maxLevel { return choice }
	}
	for _, choice := range choices {
		if choice.addon != 0 { return choice }
	}
	return CatChoice_t{ addon: 0, level: bundled, label: LevelLabel(bundled) }
}

func CategRangeMatch(plan Plan_t, categId CategId_t, minLevel, maxLevel int) bool {
	bundled := PlanBundledLevel(plan, categId)
	if bundled >= minLevel && bundled <= maxLevel { return true }

	key := PlanCateg_t{ plan: plan.planId, categ: categId }
	for _, choice := range App.lookup.planAddonChoices[key] {
		if choice.addon == 0 { continue }
		if choice.level >= minLevel && choice.level <= maxLevel { return true }
	}
	return false
}

func QuoteCommission(row QuotePlan_t, plan Plan_t, preexByCateg map[CategId_t]EuroCent_t) EuroCent_t {
	regularBase := row.planBase + preexByCateg[0]
	pvnBase := EuroCent_t(0)

	for _, addon := range row.addons {
		if !addon.priceOk || addon.base == 0 { continue }
		name := Lower(Trim(addon.categ))
		if Contains(name, `vision`) { continue }

		base := addon.base + preexByCateg[addon.categId]
		if Contains(name, `pvn`) {
			pvnBase += base
			continue
		}
		regularBase += base
	}

	upfront := Commission(regularBase, plan.comonths) + Commission(pvnBase, Months_t(200))
	commission := ApplyPercent(upfront, Percent_t(20))
	return EuroFlatFromCent(commission).ToEuroCent()
}

func VisionAmount(plan Plan_t, planBase EuroCent_t) (EuroCent_t, bool) {
	if plan.vision.percent > 0 {
		return ApplyPercent(planBase, plan.vision.percent), true
	}
	if plan.vision.euro > 0 {
		return plan.vision.euro, true
	}
	return 0, false
}

func NoClaimsAmount(planBase EuroCent_t, plan Plan_t, isChild bool) EuroCent_t {
	months := plan.nc.adult.months
	flat := plan.nc.adult.flat
	if isChild {
		months = plan.nc.child.months
		flat = plan.nc.child.flat
	}
	return flat + Commission(planBase, months)
}

func DeductAmount(plan Plan_t, isChild bool) EuroCent_t {
	if isChild { return plan.ded.child.euro }
	return plan.ded.adult.euro
}

func PlanFilters(state State_t) PlanFilters_t {
	f := PlanFilters_t{
		specref: specrefAddon,
		deductMax: int(EuroFlat_t(3000).ToEuroCent()),
		hospitalMax: 999,
		dentalMax: 999,
	}

	if v, ok := StateIntAny(state, `segment`); ok { f.segment = v }
	if v, ok := StateIntAny(state, `priorCov`, `prior`); ok { f.prior = v }
	if v, ok := StateIntAny(state, `noExam`, `exam`); ok { f.noExam = v > 0 }
	if v, ok := StateIntAny(state, `specref`, `referral`); ok { f.specref = v }
	f.vision = StateBool(state, `vision`, `glasses`)
	f.tempVisa = StateBool(state, `tempVisa`, `tempvisa`, `visa`)
	f.noPVN = StateBool(state, `noPVN`, `pvnoff`)
	f.naturalMed = StateBool(state, `naturalMed`, `natural`)

	if v, ok := StateIntAny(state, `deductibleMin`, `deductMin`); ok { f.deductMin = int(EuroFlat_t(v).ToEuroCent()) }
	if v, ok := StateIntAny(state, `deductibleMax`, `deductMax`); ok { f.deductMax = int(EuroFlat_t(v).ToEuroCent()) }
	if f.deductMin > f.deductMax { f.deductMin, f.deductMax = f.deductMax, f.deductMin }

	if v, ok := StateIntAny(state, `hospitalMin`, `minHospital`); ok { f.hospitalMin = v }
	if v, ok := StateIntAny(state, `hospitalMax`, `maxHospital`); ok { f.hospitalMax = v }
	if f.hospitalMin > f.hospitalMax { f.hospitalMin, f.hospitalMax = f.hospitalMax, f.hospitalMin }

	if v, ok := StateIntAny(state, `dentalMin`, `minDental`); ok { f.dentalMin = v }
	if v, ok := StateIntAny(state, `dentalMax`, `maxDental`); ok { f.dentalMax = v }
	if f.dentalMin > f.dentalMax { f.dentalMin, f.dentalMax = f.dentalMax, f.dentalMin }

	return f
}

func PlanFilterReasons(plan Plan_t, f PlanFilters_t) []string {
	var reasons []string

	if !f.tempVisa && plan.tempvisa { reasons = append(reasons, `Temp Visa`) }
	if f.segment > 0 && int(plan.segmask)&f.segment == 0 { reasons = append(reasons, `Segment`) }
	if f.prior < plan.priorcov { reasons = append(reasons, `Prior cover`) }
	if f.noExam && f.prior < plan.noexam { reasons = append(reasons, `No exam`) }
	if f.specref != specrefAddon && f.specref != plan.specref && plan.specref != specrefAddon { reasons = append(reasons, `Specialist`) }

	dval := int(plan.ded.adult.euro)
	if dval < f.deductMin || dval > f.deductMax { reasons = append(reasons, `Deductible`) }
	if !CategRangeMatch(plan, CategId_t(catHospital), f.hospitalMin, f.hospitalMax) { reasons = append(reasons, `Hospital`) }
	if !CategRangeMatch(plan, CategId_t(catDental), f.dentalMin, f.dentalMax) { reasons = append(reasons, `Dental`) }

	return reasons
}

func QuotePlans(state State_t) QuotePlans_t {
	buyYear, yearAge, exactAge := PlanAges(state)
	isChild := exactAge > 0 && exactAge < 21
	sickCover := StateInt(state, `sickCover`)
	filter := PlanFilters(state)
	sortBy := QuoteSortMode(StateValue(state, `sortBy`))

	out := QuotePlans_t{
		showVision: filter.vision,
		sortBy: sortBy,
	}
	for planId, plan := range App.lookup.plans.All() {
		age := yearAge
		if plan.exactAge { age = exactAge }
		suppressSurch := !plan.surcharge || filter.segment == segmentStudent
		planBase, planSurch, planOk := CatPrice(buyYear, age, planId, 0, sickCover, suppressSurch)

		reasons := PlanFilterReasons(plan, filter)
		if !planOk || planBase == 0 { reasons = append(reasons, `Zero price`) }
		if len(reasons) > 0 {
			out.filtered = append(out.filtered, QuotePlanFiltered_t{
				planId: int(plan.planId),
				label: Str(plan.provName, ` / `, plan.name),
				reasons: reasons,
			})
			continue
		}

		row := QuotePlan_t{
			planId: int(plan.planId),
			label: Str(plan.provName, ` / `, plan.name),
			deduct: DeductAmount(plan, isChild),
			noClaims: NoClaimsAmount(planBase, plan, isChild),
			planBase: planBase,
			planSurcharge: planSurch,
			planOk: planOk,
		}
		if planOk {
			row.base += planBase
			row.surcharge += planSurch
		}

		planKey := plan.planId
		for _, categ := range App.lookup.categs.All() {
			if categ.categId == 0 { continue }
			if filter.noPVN && CategIs(categ, `pvn`) { continue }

			key := PlanCateg_t{ plan: planKey, categ: categ.categId }
			choices := App.lookup.planAddonChoices[key]
			if len(choices) == 0 { continue }
			choices = EnsureFreeChoice(plan, categ.categId, choices)
			choices = EnsureOptionalChoice(categ, choices)

			choice, ok := App.lookup.planAddons[key]
			if !ok { continue }
			if IsHospDental(categ.categId) {
				minLevel, maxLevel := filter.hospitalMin, filter.hospitalMax
				if int(categ.categId) == catDental { minLevel, maxLevel = filter.dentalMin, filter.dentalMax }
				choice = PickRangeChoice(plan, categ.categId, choices, minLevel, maxLevel)
			}
			if CategIs(categ, `natural`) {
				if !filter.naturalMed {
					choice.addon = 0
					choice.label = `---`
				} else if choice.addon == 0 {
					choice.addon = PickAnyAddon(choices)
				}
			}
			if !(CategIs(categ, `natural`) && !filter.naturalMed) {
				keyName := QuotePlanCatControlName(int(plan.planId), categ.categId)
				if selected, ok := StateIntOK(state, keyName); ok {
					if CategIs(categ, `natural`) && filter.naturalMed && selected == 0 {
						// Natural Med ON should auto-pick an addon even if prior state stored 0.
					} else
					if picked, found := QuoteChoiceByAddon(choices, AddonId_t(selected)); found {
						choice = picked
					}
				}
			}

			x := QuotePlanAddon_t{
				categId: categ.categId,
				categ: categ.name,
				label: choice.label,
				addon: choice.addon,
				level: choice.level,
				hasMulti: len(choices) > 1,
				choices: choices,
			}
			if !x.hasMulti && x.addon == 0 { continue }

			if x.addon != 0 {
				x.base, x.surcharge, x.priceOk = CatPrice(buyYear, age, int(x.addon), categ.categId, sickCover, suppressSurch)
				if x.priceOk {
					row.base += x.base
					row.surcharge += x.surcharge
				}
			}

			row.addons = append(row.addons, x)
		}

		if filter.vision {
			visionBase, visionOk := VisionAmount(plan, planBase)
			if visionOk {
				row.base += visionBase
				row.addons = append(row.addons, QuotePlanAddon_t{
					categId: CategId_t(-1),
					categ: `Vision`,
					base: visionBase,
					priceOk: true,
				})
			}
		}

		row.commission = QuoteCommission(row, plan, nil)

		row.price = row.base + row.surcharge
		out.plans = append(out.plans, row)
	}

	if sortBy == sortByPrice {
		sort.Slice(out.plans, func(i, j int) bool {
			if out.plans[i].price != out.plans[j].price {
				return out.plans[i].price < out.plans[j].price
			}
			return Lower(out.plans[i].label) < Lower(out.plans[j].label)
		})
	} else {
		sort.Slice(out.plans, func(i, j int) bool {
			return Lower(out.plans[i].label) < Lower(out.plans[j].label)
		})
	}
	sort.Slice(out.filtered, func(i, j int) bool {
		return Lower(out.filtered[i].label) < Lower(out.filtered[j].label)
	})
	return out
}
