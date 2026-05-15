package system

import "testing"

func TestAlternativesForPip(t *testing.T) {
	alternatives := alternativesFor("pip")
	if len(alternatives) != 2 {
		t.Fatalf("expected 2 alternatives, got %d", len(alternatives))
	}
	if alternatives[0].Command != "pip3" {
		t.Fatalf("expected pip3 fallback, got %s", alternatives[0].Command)
	}
}
