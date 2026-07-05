package main

import "testing"

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
