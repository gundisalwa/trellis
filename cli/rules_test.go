package main

import (
	"strings"
	"testing"
)

// TestInvariantRulesCoverCatalog guards decision-0026: the always-loaded rules are
// parsed from the bundled catalog, so every assessable invariant must yield a rule.
func TestInvariantRulesCoverCatalog(t *testing.T) {
	rules := invariantRules()
	if len(rules) != 14 {
		t.Errorf("expected 14 invariant rules parsed from the catalog, got %d: %v", len(rules), sortedKeys(rules))
	}
	for _, slug := range []string{"inv-directional-flow", "floor-transparency", "floor-intent-gate", "inv-self-improvement"} {
		if rules[slug] == "" {
			t.Errorf("no rule extracted for %s", slug)
		}
	}
}

// TestInvariantPrimaryFailureCoverCatalog guards decision-0031: every assessable
// invariant must yield a primary failure example (its first `violated` case) for the
// always-loaded grounding line.
func TestInvariantPrimaryFailureCoverCatalog(t *testing.T) {
	fails := invariantPrimaryFailure()
	if len(fails) != 14 {
		t.Errorf("expected 14 primary failures parsed from the catalog, got %d: %v", len(fails), sortedKeys(fails))
	}
	// spot-check: the first violated example for directional-flow, not the second.
	if got := fails["inv-directional-flow"]; !strings.Contains(got, "still being edited") {
		t.Errorf("inv-directional-flow primary failure looks wrong: %q", got)
	}
}

// TestCatalogSlugOrder guards decision-0051 rule 4: the assembled readout
// concatenates fragments "in catalog order" — the order the entries appear in the
// bundled catalog document (structural → operating → floors), which the parser must
// preserve. The set must be exactly the assessable slugs the other parsers cover.
func TestCatalogSlugOrder(t *testing.T) {
	order := catalogSlugOrder()
	if len(order) != 14 {
		t.Fatalf("expected 14 slugs in catalog order, got %d: %v", len(order), order)
	}
	if order[0] != "inv-directional-flow" {
		t.Errorf("catalog order must open with the structural set (inv-directional-flow), got %s", order[0])
	}
	if order[12] != "floor-transparency" || order[13] != "floor-intent-gate" {
		t.Errorf("catalog order must close with the floors (floor-transparency, floor-intent-gate), got %v", order[12:])
	}
	dirs := invariantDirectives()
	for _, slug := range order {
		if dirs[slug] == "" {
			t.Errorf("catalog-order slug %s has no directive — order and directive parsers disagree", slug)
		}
	}
}

// TestSelfImprovementCarriesEntropyLean guards decision-0052: the catalog's
// inv-self-improvement entry renders the dispositional face the set already owns —
// the directive carries the pattern-introduction notice (point 1), the entry gains
// the entropy-lean signature clause (point 2) and the *(structure)* honored/violated
// pair (point 3), and the rendered rule's ✗ line extends in the same breath while
// staying one ✗ bullet, the readout format unchanged (point 4).
func TestSelfImprovementCarriesEntropyLean(t *testing.T) {
	d := invariantDirectives()["inv-self-improvement"]
	for _, want := range []string{
		"And notice the friction you are about to create",
		"migrate it, or name the exemption and ask — never resolve it silently in prose",
	} {
		if !strings.Contains(d, want) {
			t.Errorf("inv-self-improvement directive missing the decision-0052 point-1 extension %q: %q", want, d)
		}
	}
	f := invariantPrimaryFailure()["inv-self-improvement"]
	if want := "re-runs it, forever — or a new convention lands and the old stock stays loose beside it, exempted by prose nobody approved"; !strings.Contains(f, want) {
		t.Errorf("inv-self-improvement primary failure missing the decision-0052 point-4 ✗ extension: %q", f)
	}
	if frag := ruleFragment("inv-self-improvement"); strings.Count(frag, "✗") != 1 {
		t.Errorf("decision-0052 point 4: the readout format is unchanged — still exactly one ✗ bullet, got: %q", frag)
	}
	// The catalog wraps prose fields at ~100 cols, so match wrap-insensitively.
	catalog := strings.Join(strings.Fields(invariantsRef), " ")
	if !strings.Contains(catalog, "the entropy lean as proactive notice") {
		t.Error("catalog missing the decision-0052 point-2 signature clause (the entropy lean as proactive notice)")
	}
	if !strings.Contains(catalog, "migrate or exempt?") || !strings.Contains(catalog, "two conventions in one tree") {
		t.Error("catalog missing the decision-0052 point-3 *(structure)* honored/violated pair")
	}
}

// TestInvariantDirectivesCoverCatalog guards decision-0034: every invariant carries an
// imperative, host-agent-facing directive for the block — and it must not leak the
// Trellis-internal codes a host agent can't resolve.
func TestInvariantDirectivesCoverCatalog(t *testing.T) {
	dirs := invariantDirectives()
	if len(dirs) != 14 {
		t.Errorf("expected 14 directives parsed from the catalog, got %d: %v", len(dirs), sortedKeys(dirs))
	}
	for slug, d := range dirs {
		if d == "" {
			t.Errorf("%s has an empty directive", slug)
		}
		for _, code := range []string{"(A1)", "(A2)", "(A3)", "(A4)", "(B2)", "(C2)", "(D1)", "decision-0", "invariant B"} {
			if strings.Contains(d, code) {
				t.Errorf("directive for %s leaks internal code %q: %s", slug, code, d)
			}
		}
	}
}
