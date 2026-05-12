package main

import (
	"os"
	. "klpm/lib/output"
)

var AutoChoose struct {
	ready bool
	qvars QuoteVars_t
}

func DefaultVars() {
	if os.Getenv(`HOME`) != `/home/glen` { return }
	AutoChoose.qvars = QuoteDefaultQuoteVars()
	AutoChoose.ready = true

	qvars := &AutoChoose.qvars

	qvars.core.segment = 1
	AddPlan(`Inter`, `LA-VNS U`, `29S`)
	AddPlan(`HanseMerkur`, `KVS3`, `Private`)
/*
	AutoChoose.qvars.dependants = append(AutoChoose.qvars.dependants, Dependant_t{
		depId: 1,
		birth: QuoteDateAddMonths(qvars.core.buy, -12*8),
		preexByChoice: make(map[ChoiceId_t][]Preex_t),
	})
	for choiceId, choice := range AutoChoose.qvars.choices {
		choice.preex = EditQSetPreex(choice.preex, 0, editQPreexModeEur, `69`, ``, true, true, false)
		AutoChoose.qvars.choices[choiceId] = choice
	}
	qvars.lang = German;
*/
}

func AddPlan(prov, plan string, addons ...string) {
	if os.Getenv(`HOME`) != `/home/glen` { return }
	if !AutoChoose.ready {
		AutoChoose.qvars = QuoteDefaultQuoteVars()
		AutoChoose.ready = true
	}

	prov = Trim(prov)
	plan = Trim(plan)
	if prov == `` || plan == `` { return }

	providerId := ProviderID(prov)
	if providerId <= 0 { panic(Error(`provider not found: `, prov)) }

	planId := PlanID(prov, plan)
	if planId <= 0 { panic(Error(`plan not found: `, prov, ` / `, plan)) }

	if AutoChoose.qvars.choices == nil { AutoChoose.qvars.choices = make(map[ChoiceId_t]PlanQuoteInfo_t) }

	choiceId := QuoteAllocChoiceId(&AutoChoose.qvars)
	choice := PlanQuoteInfo_t{
		plan: PlanId_t(planId),
		addons: make(map[CategId_t]AddonId_t),
	}

	segment := SegmentCode(AutoChoose.qvars.core.segment)
	if segment == `` { panic(Error(`segment code not found for segment: `, AutoChoose.qvars.core.segment)) }

	for _, level := range addons {
		level = Trim(level)
		if level == `` { continue }

		levelId := LevelID(level)
		if levelId <= 0 { panic(Error(`level not found: `, level)) }

		addonId := AddonID(prov, level, segment)
		if addonId <= 0 { panic(Error(`addon not found: `, prov, ` / level `, level, ` / segment `, segment)) }

		product, ok := App.lookup.products[ProductId_t(addonId)]
		if !ok { panic(Error(`addon product not loaded: `, addonId, ` / level `, level)) }
		if product.categ <= 0 { continue }

		choice.addons[CategId_t(product.categ)] = AddonId_t(addonId)
	}

	AutoChoose.qvars.choices[choiceId] = choice
}

func ProviderID(name string) int {
	return LookupIDByName(`provider_id`, name)
}

func PlanID(provider, plan string) int {
	return LookupIDByName(`plan_id`, provider, plan)
}

func LevelID(level string) int {
	return LookupIDByName(`level_id`, level)
}

func AddonID(provider, level, segment string) int {
	return LookupIDByName(`addon_id`, provider, level, segment)
}

func SegmentCode(segment int) string {
	if segment <= 0 { return `` }
	var q Builder
	q.Add(`select code from segments where segment=`, segment)

	var out string
	row := App.DB.QueryRow(q.String()).Scan(&out)
	if row.HasError() { panic(row.Message()) }
	return Trim(out)
}

func LookupIDByName(fn string, names ...string) int {
	var q Builder
	q.Add(`select `, fn, `(`)
	for i, name := range names {
		if i > 0 { q.Add(`,`) }
		q.Add(DQ(Trim(name)))
	}
	q.Add(`)`)

	var out int
	row := App.DB.QueryRow(q.String()).Scan(&out)
	if row.HasError() { panic(row.Message()) }
	return out
}
