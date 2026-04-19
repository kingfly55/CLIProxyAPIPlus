package management

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/router-for-me/CLIProxyAPI/v6/internal/usage"
)

// knownRepoPrefixes is the set of repo prefixes used in WebCity agent slugs.
// Slugs starting with one of these are parsed as REPO-ROLE-TIER; otherwise ROLE-TIER.
var knownRepoPrefixes = map[string]bool{
	"wa": true, // workflow-automation
	"cm": true, // camoufox-manager
	"th": true, // workflow-automation-test-harness
}

// knownTiers recognises the model tier suffix.
var knownTiers = map[string]bool{
	"haiku":  true,
	"sonnet": true,
	"opus":   true,
}

// RoleRollup aggregates per-role usage with optional repo and tier breakdowns.
type RoleRollup struct {
	TotalRequests int64                   `json:"total_requests"`
	TotalTokens   int64                   `json:"total_tokens"`
	ByRepo        map[string]SlugTotals   `json:"by_repo,omitempty"`
	ByTier        map[string]SlugTotals   `json:"by_tier,omitempty"`
	BySlug        map[string]SlugTotals   `json:"by_slug,omitempty"`
}

// SlugTotals captures raw counts for a slug / repo / tier bucket.
type SlugTotals struct {
	TotalRequests int64 `json:"total_requests"`
	TotalTokens   int64 `json:"total_tokens"`
}

// GetUsageByRole returns usage aggregated by role (with optional repo/tier rollups).
//
// Expected slug format: [REPO-]ROLE-TIER (e.g. "wa-builder-sonnet", "mayor-opus",
// "th-self-check-haiku"). Slugs that don't match the format fall into "unparsed".
func (h *Handler) GetUsageByRole(c *gin.Context) {
	var snapshot usage.StatisticsSnapshot
	if h != nil && h.usageStats != nil {
		snapshot = h.usageStats.Snapshot()
	}

	roles := make(map[string]*RoleRollup)
	unparsed := make(map[string]SlugTotals)

	for _, apiSnapshot := range snapshot.APIs {
		for slug, slugSnapshot := range apiSnapshot.OriginalSlugs {
			repo, role, tier, ok := parseSlug(slug)
			if !ok {
				totals := unparsed[slug]
				totals.TotalRequests += slugSnapshot.TotalRequests
				totals.TotalTokens += slugSnapshot.TotalTokens
				unparsed[slug] = totals
				continue
			}
			rollup, exists := roles[role]
			if !exists {
				rollup = &RoleRollup{
					ByRepo: make(map[string]SlugTotals),
					ByTier: make(map[string]SlugTotals),
					BySlug: make(map[string]SlugTotals),
				}
				roles[role] = rollup
			}
			rollup.TotalRequests += slugSnapshot.TotalRequests
			rollup.TotalTokens += slugSnapshot.TotalTokens
			if repo != "" {
				addTotals(rollup.ByRepo, repo, slugSnapshot)
			}
			if tier != "" {
				addTotals(rollup.ByTier, tier, slugSnapshot)
			}
			addTotals(rollup.BySlug, slug, slugSnapshot)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"roles":    roles,
		"unparsed": unparsed,
	})
}

func addTotals(m map[string]SlugTotals, key string, s usage.SlugSnapshot) {
	totals := m[key]
	totals.TotalRequests += s.TotalRequests
	totals.TotalTokens += s.TotalTokens
	m[key] = totals
}

// parseSlug extracts (repo, role, tier) from a slug of the form [REPO-]ROLE-TIER.
// Returns ok=false if the slug doesn't end in a known tier or has fewer than 2 parts.
func parseSlug(slug string) (repo, role, tier string, ok bool) {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return "", "", "", false
	}
	parts := strings.Split(slug, "-")
	if len(parts) < 2 {
		return "", "", "", false
	}
	tier = parts[len(parts)-1]
	if !knownTiers[tier] {
		return "", "", "", false
	}
	remainder := parts[:len(parts)-1]
	if len(remainder) >= 2 && knownRepoPrefixes[remainder[0]] {
		repo = remainder[0]
		role = strings.Join(remainder[1:], "-")
	} else {
		role = strings.Join(remainder, "-")
	}
	if role == "" {
		return "", "", "", false
	}
	return repo, role, tier, true
}
