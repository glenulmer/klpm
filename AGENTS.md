# AGENTS - Guardrails

## A) Non-Negotiables

* Prefer package main for code files
* Keep `chi` + `net/http`; explicit route lines.
* Never use templates, JSON-first app flow, GOB, REST/SPA layering, ORM/DI/DTO architecture.
* Never run `gofmt`/`go fmt`/`goimports`/IDE auto-formatters, and never do style-only "idiomatic Go" rewrites.
* No go-styled detailed error plumbing

## B) Interaction Rule

* Keep feedback technical and concrete.
* Be terse, specific, and fact-based.
* No long lists of options - at most 3 clearly distinct, realistic options.
* State conclusions and recent actions first. Details afterward.
* Any path starting with an alphabetic is relative to the repository root
* If I ask for analysis, do not guess or make assumptions -- state relevant facts and conclusions succinctly.

## C) Code style
 * use all-in-one-line for if, while, switch etc if all code fits comfortably one line (less than 75 chars). eg "if count == 0 { continue }" or "if count == 0 { count++; Log(count); continue }
 * Avoid if err != nil and other code plumbing.
 * Avoid "if ok := someFunc(); ok { ... }"
 * No naked returns

## D) Code Quality Check

* Keep scope of changes as minimal as possible.
* Do not refactor or build new architecture unless I explicitly tell you to do so, and then always build a project plan in md format.
* Keep code simple and readable.
* Before adding any helper, confirm no equivalent exists in the codebase. Ask for permission to build the helper and explain why it is worth doing and whether it is duplicative.
* Do not add one-off const values - consts are for repeated use in several files. Instead use harcoded value with a comment if needed.
