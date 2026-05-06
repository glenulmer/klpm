# Responsive Migration Plan

## Objective
Move from server-selected `desktop`/`phone` rendering to a single responsive UI path while preserving existing quote/edit behavior and minimizing risk.

## Current Facts
1. Routing and static serving are in `1.main.go` with `chi` and `http.FileServer`.
2. Layout is selected server-side per request through `RequestLayout(req)` and `SessionDeviceMode(req)`.
3. `page.1.quote.go` and `page.2.editq.go` pick layout-specific CSS files.
4. Quote has separate view trees in `z.render.desktop.go` and `z.render.phone.go`.
5. EditQ has separate entry points but nearly identical body structure (`page.2.editq.render.desktop.go`, `page.2.editq.render.phone.go`).
6. Rewrite endpoints (`page.1.quote.post.go`, `page.2.editq.post.go`) return HTML fragments and depend on stable IDs and controls.

## Migration Strategy
Use staged PRs to avoid a big-bang switch.

## PR1: Unify EditQ to one responsive path
### Scope
1. Keep current behavior and endpoints.
2. Remove layout split only for EditQ rendering and CSS loading.
3. Keep all form names, ids, and rewrite targets unchanged.

### File changes
1. `page.2.editq.go`
2. `page.2.editq.post.go`
3. `page.2.editq.render.go`
4. `page.2.editq.render.desktop.go`
5. `page.2.editq.render.phone.go`
6. `static/css/page.2.editq.desktop.css`
7. `static/css/page.2.editq.phone.css`
8. `static/js/page.2.editq.js` (only if selector coupling appears during testing)

### Implementation steps
1. Introduce one EditQ CSS entrypoint: `static/css/page.2.editq.css`.
2. Point `Page2EditQ` to this single CSS.
3. Collapse `EditQBodyView` to one branch returning a unified body class.
4. Keep `EditQFormBody` ID unchanged for rewrite compatibility.
5. Build responsive breakpoints inside `page.2.editq.css` (mobile-first, desktop overrides).
6. Keep visual parity with existing desktop/phone at their typical widths.

### Acceptance checks
1. `/quote-review` GET renders same controls and section order.
2. `/quote-review-change` rewrite updates still apply and preserve caret/focus behavior.
3. Dependants add/remove works.
4. Download flow (`DownloadExcel`) still works.
5. Visual checks at 320, 390, 768, 1024, 1440 widths.

### Rollback
1. Revert CSS path switch in `page.2.editq.go`.
2. Restore old `EditQBodyView` branch if needed.

## PR2: Unify Quote to one responsive path
### Scope
1. Keep quote logic, selection logic, and rewrite endpoints unchanged.
2. Replace desktop/phone rendering split with one responsive markup.
3. Replace desktop/phone CSS split with one responsive stylesheet.

### File changes
1. `page.1.quote.go`
2. `page.1.quote.post.go`
3. `z.render.desktop.go`
4. `z.render.phone.go`
5. `z.quote.plan.render.go`
6. `static/css/page.1.quote.desktop.css`
7. `static/css/page.1.quote.phone.css`
8. `static/js/page.1.quote.js`

### Implementation steps
1. Add a single CSS entrypoint: `static/css/page.1.quote.css`.
2. Point `QuoteCSSPath` to one stylesheet.
3. Create a unified quote page view preserving IDs:
   1. `QuoteFormBody`
   2. `QuotePlans`
   3. `QuoteInfoCard`
   4. `QuoteSelectedCard`
   5. `QuotePhoneStickyAnchor` (can remain for compatibility during transition)
4. Preserve all control `name` attributes and button name patterns (`seladd-*`, `seldel-*`, `selcat-*`, `quoteReset`).
5. Convert fixed desktop width constraints to responsive breakpoints.
6. Keep desktop table and mobile card behavior available through CSS breakpoint visibility if full structural unification is too risky in one PR.

### Acceptance checks
1. `/` GET renders and all inputs trigger `/quote-info-change` updates.
2. Plan add/remove and selected addons still work.
3. Sort select still rewrites results.
4. Sticky behavior remains correct on narrow viewports.
5. Visual checks at 320, 390, 768, 1024, 1280, 1440 widths.

### Rollback
1. Revert CSS path in `QuoteCSSPath`.
2. Restore old desktop/phone branch in `QuotePageView` and `RewriteQuotePage`.

## PR3: Remove server layout switching
### Scope
Clean up layout/session mechanisms once both pages are responsive.

### File changes
1. `z.layout.detect.go`
2. `z.sessions.go`
3. `z.state.go`
4. `0.login.go`
5. `page.1.quote.go`
6. `page.2.editq.go`
7. `page.1.quote.post.go`
8. `page.2.editq.post.go`

### Implementation steps
1. Remove `layout` query processing in middleware.
2. Remove `device` and `deviceConfirmed` session fields.
3. Remove `DeviceConfirmHeadScript` injection on signin.
4. Keep session and auth logic intact.
5. Remove dead helper funcs tied to layout split.

### Acceptance checks
1. Login/logout/session flow unchanged.
2. No route depends on `?layout=`.
3. No server-selected layout path remains.

### Rollback
1. Reintroduce `layout` query handling and session device fields.

## PR4: Delete obsolete files and tighten checks
### Scope
Delete old split assets/files and enforce no reintroduction.

### File changes
1. Remove obsolete CSS and render split files after PR1/PR2/PR3 are stable.
2. Update any script/check references that still mention old files.

### Acceptance checks
1. Build and runtime start cleanly.
2. All routes load with expected CSS/JS.

## Testing Matrix (run each PR)
1. Auth: signin, signout, session persistence.
2. Quote: core edits, plan selection, selected plan updates, reset.
3. EditQ: dependant edits, preex changes, review table, download.
4. Rewrite reliability: repeated fast input changes and no broken DOM targets.
5. Responsive visual checks at defined breakpoints.

## Sequencing
1. Merge PR1.
2. Merge PR2.
3. Merge PR3.
4. Merge PR4.

## Notes
1. Do not change backend quote math or DB access in this migration.
2. Keep rewrite target IDs stable until PR3 cleanup.
3. Keep scope minimal per PR and avoid unrelated refactors.
