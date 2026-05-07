package main

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/css"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
)

type scenario_t struct {
	name   string
	path   string
	width  int64
	height int64
}

type rule_t struct {
	source    string
	start     int
	end       int
	line      int
	selector  string
	snippet   string
	used      bool
	byScenario map[string]bool
}

func stripQuery(u string) string {
	p, e := url.Parse(u)
	if e != nil { return u }
	p.RawQuery = ``
	p.Fragment = ``
	return p.String()
}

func lineAt(text string, offset int) int {
	if offset <= 0 { return 1 }
	if offset > len(text) { offset = len(text) }
	return strings.Count(text[:offset], "\n") + 1
}

func selectorFromSnippet(snippet string) string {
	x := strings.TrimSpace(snippet)
	if x == `` { return `` }
	i := strings.Index(x, "{")
	if i < 0 { return x }
	return strings.TrimSpace(x[:i])
}

func oneLine(s string) string {
	s = strings.TrimSpace(strings.ReplaceAll(s, "\n", " "))
	s = strings.Join(strings.Fields(s), " ")
	if len(s) > 140 { return s[:140] + ` ...` }
	return s
}

func main() {
	baseURL := flag.String(`base`, `http://127.0.0.1:3001`, `base URL to audit`)
	cssNeedle := flag.String(`css`, `/static/css/responsive.css`, `css URL (path or full URL) to report`)
	browserPath := flag.String(`browser`, ``, `optional browser executable path`)
	flag.Parse()

	scenarios := []scenario_t{
		{ name:`desktop`, path:`/`, width:1366, height:900 },
		{ name:`phone`, path:`/`, width:390, height:844 },
	}

	allocOpts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag(`headless`, true),
		chromedp.Flag(`disable-gpu`, true),
		chromedp.Flag(`no-sandbox`, true),
		chromedp.Flag(`hide-scrollbars`, true),
		chromedp.Flag(`disable-crash-reporter`, true),
		chromedp.Flag(`disable-breakpad`, true),
		chromedp.Flag(`no-first-run`, true),
		chromedp.Flag(`no-default-browser-check`, true),
		chromedp.UserDataDir(`/tmp/chromedp-csscov-profile`),
	)
	if *browserPath != `` { allocOpts = append(allocOpts, chromedp.ExecPath(*browserPath)) }

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), allocOpts...)
	defer cancelAlloc()

	ctx, cancelCtx := chromedp.NewContext(allocCtx)
	defer cancelCtx()

	byID := map[cdp.StyleSheetID]*css.StyleSheetHeader{}
	var byIDLock sync.Mutex
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		x, ok := ev.(*css.EventStyleSheetAdded)
		if !ok || x == nil { return }
		byIDLock.Lock()
		byID[x.Header.StyleSheetID] = x.Header
		byIDLock.Unlock()
	})

	rules := map[string]*rule_t{}

	for _, sc := range scenarios {
		targetURL := strings.TrimRight(*baseURL, `/`) + sc.path
		var usage []*css.RuleUsage

		e := chromedp.Run(ctx,
			chromedp.ActionFunc(func(ctx context.Context) error { return emulation.SetDeviceMetricsOverride(sc.width, sc.height, 1, false).Do(ctx) }),
			chromedp.ActionFunc(func(ctx context.Context) error { return css.Enable().Do(ctx) }),
			chromedp.ActionFunc(func(ctx context.Context) error { return css.StartRuleUsageTracking().Do(ctx) }),
			chromedp.Navigate(targetURL),
			chromedp.WaitReady(`body`, chromedp.ByQuery),
			chromedp.Sleep(450*time.Millisecond),
			chromedp.Evaluate(`(() => {
				const toggle = (sel, count) => {
					const xs = Array.from(document.querySelectorAll(sel));
					for (let i = 0; i < xs.length && i < count; i++) xs[i].click();
				};
				window.scrollTo(0, document.body.scrollHeight);
				window.scrollTo(0, 0);
				toggle('#QuoteInfoCard > summary', 1);
				toggle('#QuoteInfoCard > summary', 1);
				toggle('#QuoteSelectedCard > summary', 1);
				toggle('#QuoteSelectedCard > summary', 1);
				toggle('details.quote-plan-addon-details > summary', 4);
				const sel = Array.from(document.querySelectorAll('#QuoteForm select[name], #QuotePlans select[name]'));
				for (let i = 0; i < sel.length && i < 4; i++) {
					const x = sel[i];
					if (!(x instanceof HTMLSelectElement) || x.options.length < 2) continue;
					x.selectedIndex = Math.min(1, x.options.length - 1);
					x.dispatchEvent(new Event('input', { bubbles:true }));
					x.dispatchEvent(new Event('change', { bubbles:true }));
				}
			})()`, nil),
			chromedp.Sleep(900*time.Millisecond),
			chromedp.ActionFunc(func(ctx context.Context) error {
				var e error
				usage, e = css.StopRuleUsageTracking().Do(ctx)
				return e
			}),
		)
		if e != nil {
			fmt.Fprintf(os.Stderr, "scenario %s failed: %v\n", sc.name, e)
			os.Exit(1)
		}

		textByID := map[cdp.StyleSheetID]string{}

		for _, u := range usage {
			if u == nil { continue }
			byIDLock.Lock()
			header, ok := byID[u.StyleSheetID]
			byIDLock.Unlock()
			if !ok { continue }
			src := header.SourceURL
			if src == `` { continue }
			src0 := stripQuery(src)
			needle := *cssNeedle
			if strings.HasPrefix(needle, `http://`) || strings.HasPrefix(needle, `https://`) {
				needle = stripQuery(needle)
			}
			if !strings.Contains(src0, needle) { continue }

			text, ok := textByID[u.StyleSheetID]
			if !ok {
				var e error
				text, e = css.GetStyleSheetText(u.StyleSheetID).Do(ctx)
				if e != nil { continue }
				textByID[u.StyleSheetID] = text
			}

			start, end := int(u.StartOffset), int(u.EndOffset)
			if start < 0 || end <= start || end > len(text) { continue }
			snippet := text[start:end]
			key := fmt.Sprintf(`%s:%d:%d`, src0, start, end)
			r, ok := rules[key]
			if !ok {
				r = &rule_t{
					source: src0,
					start: start,
					end: end,
					line: lineAt(text, start),
					selector: selectorFromSnippet(snippet),
					snippet: oneLine(snippet),
					byScenario: map[string]bool{},
				}
				rules[key] = r
			}
			r.byScenario[sc.name] = bool(u.Used)
			if u.Used { r.used = true }
		}
	}

	var all []*rule_t
	for _, r := range rules { all = append(all, r) }
	sort.Slice(all, func(i, j int) bool {
		if all[i].source != all[j].source { return all[i].source < all[j].source }
		if all[i].line != all[j].line { return all[i].line < all[j].line }
		return all[i].start < all[j].start
	})

	used, unused := 0, 0
	for _, r := range all {
		if r.used { used++; continue }
		unused++
	}

	fmt.Printf("CSS coverage report\n")
	fmt.Printf("Target CSS: %s\n", *cssNeedle)
	fmt.Printf("Scenarios: desktop(1366x900), phone(390x844)\n")
	fmt.Printf("Rules seen: %d | used: %d | unused: %d\n\n", len(all), used, unused)

	if len(all) == 0 {
		fmt.Println("No matching CSS rules tracked. Check -base, -css, and that the page loaded.")
		return
	}

	fmt.Println("Unused rules:")
	for _, r := range all {
		if r.used { continue }
		sel := r.selector
		if sel == `` { sel = r.snippet }
		fmt.Printf("- %s:%d :: %s\n", r.source, r.line, oneLine(sel))
	}
}
