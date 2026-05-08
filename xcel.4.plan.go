package main

import (
	_ "image/jpeg"
	"os"

	. "klpm/lib/date"
	. "klpm/lib/dec2"
	. "klpm/lib/output"
)

var planColumns = []string{ `D`, `F`, `H`, `J`, `L` }

func (xl *Excel_t)WritePlanInfo(lastBenefitRow int, col string, item QuoteSelectedItem_t, row QuotePlan_t, plan Plan_t) {
	xl.writePlanLogoCell(col, plan)
	xl.writePlanTopNoteCell(col, plan)
	xl.writePlanLabelCell(col, plan)
	xl.writePlanRefCell(col, row, plan)
	xl.writePlanMainCostCells(col, item, row)
	xl.writePlanDependantCells(col, item, row)
	xl.writePlanTotalsCells(col, item, row)
	xl.writePlanBenefitValues(col, item, row, plan)
	xl.writePlanTipValues(lastBenefitRow, col, plan)
}

func (xl *Excel_t)writePlanLogoCell(col string, plan Plan_t) {
	logo := Str(`assets/logos/`, Lower(plan.provName), `.jpg`)
	if _, e := os.Stat(logo); e != nil {
		Log(`excel logo not found for provider: `, plan.provName)
		return
	}
	if e := xl.AddPicture(quoteSheet, Str(col, 2), logo, nil); e != nil { Log(`excel add logo failed for provider: `, plan.provName, ` / `, e) }
}

func (xl *Excel_t)writePlanTopNoteCell(col string, plan Plan_t) {
	lang := xl.qvars.lang
	if lang <= 0 { lang = English }
	note := App.lookup.planNotes[TopNote_t{ plan:plan.planId, lang:lang }]
	if note == `` && lang != English { note = App.lookup.planNotes[TopNote_t{ plan:plan.planId, lang:English }] }
	if Trim(note) == `` { return }
	_ = xl.SetCellValue(quoteSheet, Str(col, 3), note)
	noteCell := Str(col, 3)
	if style := xl.styles[`TopNote`]; style != 0 { _ = xl.SetCellStyle(quoteSheet, noteCell, noteCell, style) }
}

func (xl *Excel_t)writePlanLabelCell(col string, plan Plan_t) {
	_ = xl.SetCellValue(quoteSheet, Str(col, 4), Str(plan.provName, ` / `, plan.name))
}

func (xl *Excel_t)writePlanRefCell(col string, row QuotePlan_t, plan Plan_t) {
	cell := Str(col, 5)
	_ = xl.SetCellValue(quoteSheet, cell, Str(`Ref `, plan.planId, `/`, xl.commRefText(row.commission)))
	if style := xl.styles[`Commission`]; style != 0 { _ = xl.SetCellStyle(quoteSheet, cell, cell, style) }
}

func (xl *Excel_t)writePlanMainCostCells(col string, item QuoteSelectedItem_t, row QuotePlan_t) {
	preexByItem, _ := EditQReviewPreexByItemCateg(xl.qvars)
	preex := preexByItem[item.itemId]

	pvn := row.price - row.price
	if addon, ok := QuotePlanAddonByTag(row, `pvn`); ok && addon.priceOk { pvn = addon.base + addon.surcharge }

	sick := row.price - row.price
	after := ``
	if addon, ok := QuotePlanAddonByTag(row, `sick`); ok {
		if addon.priceOk { sick = addon.base + addon.surcharge }
		after = xl.sickAfterText(addon)
	}

	total := row.price + preex
	hic := total - (pvn + sick + preex)

	mainCell := Str(col, 6)
	if preex > 0 {
		_ = xl.SetCellValue(quoteSheet, mainCell, Str(hic.OutEuro(), ` + `, preex.OutEuro()))
	} else {
		_ = xl.SetCellValue(quoteSheet, mainCell, hic.OutEuro())
	}
	_ = xl.SetCellValue(quoteSheet, Str(col, 7), pvn.OutEuro())
	_ = xl.SetCellValue(quoteSheet, Str(col, 8), sick.OutEuro())
	_ = xl.SetCellValue(quoteSheet, Str(col, 9), after)
}

func (xl *Excel_t)writePlanDependantCells(col string, item QuoteSelectedItem_t, row QuotePlan_t) {
	for i, dep := range xl.qvars.dependants {
		if i >= 10 { break }
		depId := dep.depId
		if depId <= 0 { depId = i + 1 }
		birth := ``
		if Valid(dep.birth) { birth = dep.birth.Format(`yyyy-mm-dd`) }
		charges := EditQDependantCharges(xl.qvars, EditQDep_t{
			depId: depId,
			name: dep.name,
			birth: birth,
			vision: dep.vision,
		})
		base := row.price - row.price
		preex := row.price - row.price
		found := false
		for _, x := range charges {
			if x.itemId != item.itemId { continue }
			if !found { base = x.planPrice; found = true }
			preex += x.applied
		}
		if !found { continue }
		_ = xl.SetCellValue(quoteSheet, Str(col, 11+i), (base+preex).OutEuro())
	}
}

func (xl *Excel_t)writePlanTotalsCells(col string, item QuoteSelectedItem_t, row QuotePlan_t) {
	preexByItem, preexByItemCateg := EditQReviewPreexByItemCateg(xl.qvars)
	total := row.price + preexByItem[item.itemId]
	_ = xl.SetCellValue(quoteSheet, Str(col, 21), total.OutEuro())

	youpay := total
	if xl.qvars.core.segment == Employee {
		pvn := row.price - row.price
		pvnCateg := CategId_t(0)
		if addon, ok := QuotePlanAddonByTag(row, `pvn`); ok && addon.priceOk {
			pvn = addon.base + addon.surcharge
			pvnCateg = addon.categId
		}
		if pvnCateg != 0 {
			key := Str(item.itemId, `:`, pvnCateg)
			pvn += preexByItemCateg[key]
		}

		wo := total - pvn
		companyPays := wo / 2

		year := CurrentDBDate().Year()
		if Valid(xl.qvars.core.buy) { year = xl.qvars.core.buy.Year() }
		if y, ok := App.lookup.years.byId[year]; ok {
			maxshare := y.maxshare.ToEuroCent()
			if companyPays > maxshare { companyPays = maxshare }
		}

		companyPays += pvn / 2
		youpay -= companyPays
	}

	for i, dep := range xl.qvars.dependants {
		if i >= 10 { break }
		depId := dep.depId
		if depId <= 0 { depId = i + 1 }
		birth := ``
		if Valid(dep.birth) { birth = dep.birth.Format(`yyyy-mm-dd`) }
		charges := EditQDependantCharges(xl.qvars, EditQDep_t{
			depId: depId,
			name: dep.name,
			birth: birth,
			vision: dep.vision,
		})
		base := row.price - row.price
		preex := row.price - row.price
		found := false
		for _, x := range charges {
			if x.itemId != item.itemId { continue }
			if !found { base = x.planPrice; found = true }
			preex += x.applied
		}
		if !found { continue }
		youpay += base + preex
	}
	_ = xl.SetCellValue(quoteSheet, Str(col, 22), youpay.OutEuro())
}

func (xl *Excel_t)writePlanBenefitValues(col string, item QuoteSelectedItem_t, row QuotePlan_t, plan Plan_t) {
	lang := int(xl.qvars.lang)
	if lang <= 0 { lang = int(English) }
	benefitRow := 24

	for _, sec := range App.lookup.benSecs.All() {
		if sec.lang != lang { continue }
		items := xl.benefitItems(sec.section, lang, xl.qvars.slim)
		for _, ben := range items {
			offer := ``
			valueStyle := xl.sectionValueStyle(sec.section)

			switch ben.benefit {
			case 1:
				offer = xl.deductYearText(row.deduct)
				if style := xl.styles[`Spec1Value`]; style != 0 { valueStyle = style }
			case 2:
				offer = xl.noClaimsText(row.noClaims)
			default:
				offer = App.lookup.bensByFamily[BenFamily_t{ benefit:ben.benefit, family:int(plan.familyId), lang:lang }]
				for _, addon := range row.addons {
					if addon.addon == 0 { continue }
					x, ok := App.lookup.bensByAddon[BenAddon_t{ benefit:ben.benefit, addon:int(addon.addon), lang:lang }]
					if !ok { continue }
					offer = x
					break
				}
				if offer == `` && lang != int(English) {
					offer = App.lookup.bensByFamily[BenFamily_t{ benefit:ben.benefit, family:int(plan.familyId), lang:int(English) }]
					for _, addon := range row.addons {
						if addon.addon == 0 { continue }
						x, ok := App.lookup.bensByAddon[BenAddon_t{ benefit:ben.benefit, addon:int(addon.addon), lang:int(English) }]
						if !ok { continue }
						offer = x
						break
					}
				}
			}

			cell := Str(col, benefitRow)
			_ = xl.SetCellValue(quoteSheet, cell, offer)
			if valueStyle != 0 { _ = xl.SetCellStyle(quoteSheet, cell, cell, valueStyle) }
			benefitRow++
		}
		benefitRow++
	}
}

func (xl *Excel_t)writePlanTipValues(lastBenefitRow int, col string, plan Plan_t) {
	lang := xl.qvars.lang
	if lang <= 0 { lang = English }
	tips := App.lookup.familyTips[FamilyTip_t{ family:plan.familyId, lang:lang }]
	if len(tips) == 0 && lang != English { tips = App.lookup.familyTips[FamilyTip_t{ family:plan.familyId, lang:English }] }
	if len(tips) == 0 { return }

	row := lastBenefitRow + 2
	style := xl.styles[`TipText`]
	for i, tip := range tips {
		cell := Str(col, row+i)
		_ = xl.SetCellValue(quoteSheet, cell, tip)
		if style != 0 { _ = xl.SetCellStyle(quoteSheet, cell, cell, style) }
	}
}

func (xl *Excel_t)sectionValueStyle(section int) int {
	switch section {
	case 0:
		return xl.styles[`Sec0Value`]
	case 1:
		return xl.styles[`Sec1Value`]
	case 2:
		return xl.styles[`Sec2Value`]
	case 3:
		return xl.styles[`Sec3Value`]
	case 4:
		return xl.styles[`Sec4Value`]
	}
	return 0
}

func (xl *Excel_t)deductYearText(amount EuroCent_t) string {
	whole := EuroFlatFromCent(amount)
	switch xl.qvars.lang {
	case German:
		return Str(whole, ` € / Jahr`)
	default:
		return Str(whole, ` € / year`)
	}
}

func (xl *Excel_t)noClaimsText(amount EuroCent_t) string {
	whole := EuroFlatFromCent(amount)
	switch xl.qvars.lang {
	case German:
		if amount > 0 { return Str(`Rückerstattung von `, whole, ` € / Jahr möglich`) }
		return `Keine Erstattung`
	default:
		if amount > 0 { return Str(`Refund of `, whole, ` € / year possible`) }
		return `No refund is available`
	}
}

func (xl *Excel_t)commRefText(amount EuroCent_t) string {
	return Str(EuroFlatFromCent(amount))
}

func (xl *Excel_t)sickAfterText(addon QuotePlanAddon_t) string {
	day := 0
	if addon.level > 0 {
		if addon.level%2 == 1 { day = 43 } else { day = 29 }
	}
	if day == 0 {
		code := Upper(Trim(QuoteAddonPickText(addon)))
		if Contains(code, `43`) { day = 43 }
		if Contains(code, `29`) { day = 29 }
	}
	if day == 0 { day = 43 }

	daily := (int64(xl.qvars.core.sickCover) / 4500) * 10
	if daily <= 0 {
		switch xl.qvars.lang {
		case German:
			return `Nicht gewählt`
		default:
			return `Not selected`
		}
	}

	suffix := `th`
	switch day {
	case 1, 21, 31, 41:
		suffix = `st`
	case 2, 22, 32, 42:
		suffix = `nd`
	case 3, 23, 33, 43:
		suffix = `rd`
	}
	switch xl.qvars.lang {
	case German:
		return Str(daily, ` €/Tag ab dem `, day, `. Tag`)
	default:
		return Str(daily, ` €/day as of `, day, suffix, ` day`)
	}
}
