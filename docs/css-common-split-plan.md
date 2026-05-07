# CSS Common Split Plan

## Goal
Split current page CSS into:
- `static/css/common.css` for shared identical structures
- `static/css/choose.css` for page 1 specific styles
- `static/css/review.css` for page 2 specific styles
- keep `static/css/date-ctrl.css` independent

## Scope
1. Rename:
   - `responsive.css` -> `choose.css`
   - `responsive2.css` -> `review.css`
2. Extract truly shared base blocks into `common.css`.
3. Remove extracted blocks from `choose.css` and `review.css`.
4. Update handlers to load:
   - page 1: `common.css` + `choose.css` + `date-ctrl.css`
   - page 2: `common.css` + `review.css`
5. Verify no stale references remain.

## Risk Control
- Only extract blocks that are text-identical or selector-group equivalent.
- Keep page-specific layout and media-query sections untouched.
- Verify with grep-based reference checks after refactor.
