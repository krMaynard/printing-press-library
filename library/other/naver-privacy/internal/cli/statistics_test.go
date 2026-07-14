// Copyright 2026 Kieran Maynard and contributors. Licensed under Apache-2.0. See LICENSE.
// Hand-authored tests for the novel transparency commands' pure helpers.

package cli

import "testing"

func TestNaverPeriodLabel(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"상반기", "1H"},
		{"하반기", "2H"},
		{" 상반기 ", "1H"},
		{"H1", "H1"},
	}
	for _, tc := range cases {
		if got := naverPeriodLabel(tc.in); got != tc.want {
			t.Errorf("naverPeriodLabel(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestStripHTML(t *testing.T) {
	in := `<ul><li><p class="p_dsc">압수영장은 4,355건 중 3,028건</p></li><li><p>두 번째 &amp; 항목<br/>줄바꿈</p></li></ul>`
	got := stripHTML(in)
	want := "압수영장은 4,355건 중 3,028건\n\n두 번째 & 항목\n줄바꿈"
	if got != want {
		t.Errorf("stripHTML = %q, want %q", got, want)
	}
}

func TestNaverCategoriesCoverAllColumnGroups(t *testing.T) {
	s := naverStatistic{
		WarrantRequestCount:              "1",
		CommDataRequestCount:             "2",
		CommRestrictionRequestCount:      "3",
		CommConfirmationDataRequestCount: "4",
	}
	seen := map[string]bool{}
	for _, cat := range naverCategories {
		requests, _, _, _, _ := cat.pick(s)
		if requests == "" {
			t.Errorf("category %s picked empty requests column", cat.slug)
		}
		seen[requests] = true
	}
	if len(seen) != 4 {
		t.Errorf("expected the 4 categories to map to 4 distinct column groups, got %d", len(seen))
	}
}
