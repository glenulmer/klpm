package main

import . "klpm/lib/htmlHelper"

func QuotePhoneFormBodyView(vars QuoteVars_t) Elem_t {
	return Div().Id(`QuoteFormBody`).Wrap(
		Div().Id(`QuotePhoneStickyAnchor`).Class(`quote-phone-sticky-anchor`),
		Elem(`details`).Id(`QuoteInfoCard`).Class(`quote-card`, `quote-phone-card`, `quote-phone-fold`, `quote-phone-info-fold`).KV(`open`, `open`).Wrap(
			Elem(`summary`).Class(`quote-card-title`, `quote-phone-fold-title`, `quote-phone-selected-title`).Wrap(
				Span(`Quote Info`),
				Div().Class(`quote-phone-actions`).Wrap(
					Elem(`button`).
						Type(`submit`).
						KV(`formaction`, `/signout`).
						KV(`formmethod`, `get`).
						KV(`formnovalidate`, `formnovalidate`).
						Class(`quote-edit-quote-btn`, `quote-phone-selected-edit-btn`).
						Text(`Logout`),
					QuoteResetButton(`quote-phone-selected-edit-btn`),
				),
			),
			Div().Class(`quote-info-layout`).Wrap(
				Div().Class(`quote-info-body`, `quote-info-grid-desktop`).Wrap(
					Div().Class(`quote-desk-row`, `quote-desk-row-top`).Wrap(
						QuoteNamedControlOnlySpanView(layoutDesktop, `clientName`, vars, 0, `quote-desk-no-label`, `quote-desk-name`),
						QuoteNamedControlOnlySpanView(layoutDesktop, `segment`, vars, 0, `quote-desk-no-label`, `quote-desk-segment`),
					),
					Div().Class(`quote-desk-row`, `quote-desk-row-mid`).Wrap(
						QuoteNamedControlSpanView(layoutDesktop, `birth`, vars, 0, `quote-desk-compact`),
						QuoteNamedControlSpanView(layoutDesktop, `buy`, vars, 0, `quote-desk-compact`),
						QuoteNamedControlSpanView(layoutDesktop, `sickCover`, vars, 0, `quote-desk-compact`, `quote-desk-right`),
					),
					Div().Class(`quote-desk-row`, `quote-desk-row-mid`).Wrap(
						QuoteNamedControlSpanView(layoutDesktop, `priorCov`, vars, 0, `quote-desk-compact`),
						QuoteNamedControlSpanView(layoutDesktop, `exam`, vars, 0, `quote-desk-compact`),
						QuoteNamedControlSpanView(layoutDesktop, `specref`, vars, 0, `quote-desk-compact`),
					),
					Div().Class(`quote-desk-row`, `quote-desk-row-flags`).Wrap(
						Div().Class(`quote-desk-flags-wrap`).Wrap(
							Div().Class(`quote-desk-flags`).Wrap(
								QuoteNamedControlSpanView(layoutDesktop, `vision`, vars, 1, `quote-desk-flag`),
								QuoteNamedControlSpanView(layoutDesktop, `tempVisa`, vars, 1, `quote-desk-flag`),
								QuoteNamedControlSpanView(layoutDesktop, `noPVN`, vars, 1, `quote-desk-flag`),
								QuoteNamedControlSpanView(layoutDesktop, `naturalMed`, vars, 1, `quote-desk-flag`),
							),
						),
					),
					Div().Class(`quote-desk-row`, `quote-desk-row-bottom`).Wrap(
						QuoteNamedControlLabelSpanView(layoutDesktop, `deductibleMin`, `Deductible Min`, vars, 0, `quote-desk-compact`, `quote-desk-right`),
						QuoteNamedControlLabelSpanView(layoutDesktop, `hospitalMin`, `Hospital Min`, vars, 0, `quote-desk-compact`),
						QuoteNamedControlLabelSpanView(layoutDesktop, `dentalMin`, `Dental Min`, vars, 0, `quote-desk-compact`),
					),
					Div().Class(`quote-desk-row`, `quote-desk-row-bottom`).Wrap(
						QuoteNamedControlOnlySpanView(layoutDesktop, `deductibleMax`, vars, 0, `quote-desk-no-label`, `quote-desk-compact`, `quote-desk-right`),
						QuoteNamedControlOnlySpanView(layoutDesktop, `hospitalMax`, vars, 0, `quote-desk-no-label`, `quote-desk-compact`),
						QuoteNamedControlOnlySpanView(layoutDesktop, `dentalMax`, vars, 0, `quote-desk-no-label`, `quote-desk-compact`),
					),
				),
			),
		),
		QuotePhoneSelectedPlansBox(vars),
	)
}

func QuotePhoneFormView(vars QuoteVars_t) Elem_t {
	return Elem(`form`).
		Id(`QuoteForm`).
		Class(`quote-form`, `quote-form-phone`).
		KV(`method`, `post`).
		KV(`action`, `/quote-info-change`).
		Wrap(QuotePhoneFormBodyView(vars))
}

func QuotePhonePlansView(data QuotePlans_t) Elem_t {
	var plans []Elem_t
	for _, x := range data.plans { plans = append(plans, QuotePlanCardView(x)) }
	return Div().Id(`QuotePlans`).Class(`quote-plan-results`, `quote-phone-results`).Wrap(
		Elem(`details`).Class(`quote-phone-fold`, `quote-phone-plans-fold`).KV(`open`, `open`).Wrap(
			Elem(`summary`).Class(`quote-plan-toolbar`, `quote-phone-plan-toolbar`, `quote-phone-fold-title`).Wrap(
				Div(`Plans (` , len(data.plans), `)`).Class(`quote-card-title`),
				Div().Class(`quote-plan-sort`).Wrap(
					QuoteSortSelectView(data.sortBy),
				),
			),
			Div().Class(`quote-phone-plans-body`).Wrap(
				Div().Class(`quote-plan-list`, `quote-plan-list-phone`).Wrap(plans),
			),
		),
		QuoteFilteredPlansBox(data.filtered),
	)
}

func QuotePhonePageView(vars QuoteVars_t, plans QuotePlans_t) Elem_t {
	return Elem(`main`).Class(`quote-page`, `quote-page-phone`).Wrap(
		QuotePhoneFormView(vars),
		QuotePhonePlansView(plans),
	)
}
