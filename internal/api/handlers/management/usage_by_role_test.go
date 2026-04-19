package management

import "testing"

func TestParseSlug(t *testing.T) {
	cases := []struct {
		name     string
		slug     string
		wantRepo string
		wantRole string
		wantTier string
		wantOK   bool
	}{
		{"repo+role+tier", "wa-builder-sonnet", "wa", "builder", "sonnet", true},
		{"repo+hyphenated-role+tier", "th-self-check-haiku", "th", "self-check", "haiku", true},
		{"role+tier no repo", "mayor-opus", "", "mayor", "opus", true},
		{"hyphenated role no repo", "self-check-sonnet", "", "self-check", "sonnet", true},
		{"cm prefix", "cm-auditor-haiku", "cm", "auditor", "haiku", true},
		{"unknown prefix treated as role", "unknown-auditor-opus", "", "unknown-auditor", "opus", true},
		{"unknown tier rejected", "wa-builder-ultra", "", "", "", false},
		{"empty", "", "", "", "", false},
		{"single segment", "sonnet", "", "", "", false},
		{"only tier with hyphen prefix", "-sonnet", "", "", "", false},
		{"gpt-5.4 native (non-role) rejected", "gpt-5.4-mini(high)", "", "", "", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo, role, tier, ok := parseSlug(tc.slug)
			if ok != tc.wantOK {
				t.Fatalf("ok=%v, want %v (repo=%q role=%q tier=%q)", ok, tc.wantOK, repo, role, tier)
			}
			if !ok {
				return
			}
			if repo != tc.wantRepo {
				t.Errorf("repo=%q, want %q", repo, tc.wantRepo)
			}
			if role != tc.wantRole {
				t.Errorf("role=%q, want %q", role, tc.wantRole)
			}
			if tier != tc.wantTier {
				t.Errorf("tier=%q, want %q", tier, tc.wantTier)
			}
		})
	}
}
