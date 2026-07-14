// Copyright 2026 Kieran Maynard and contributors. Licensed under Apache-2.0. See LICENSE.
// Hand-authored support for the novel transparency commands (statistics,
// summaries, whitepapers).

package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"naver-privacy-pp-cli/internal/client"

	nethtml "html"
)

const (
	naverStatisticsPageCode = "TRANSPARENCY_REPORT_STATISTICS"
	naverReportPageCode     = "NAVER_REPORT"
)

// naverStatistic is one half-year entry of the transparency series as the CMS
// delivers it: 24 wide per-category columns plus the narrative HTML.
type naverStatistic struct {
	ID              string `json:"id"`
	Year            int    `json:"year"`
	Period          string `json:"period"`
	GuideEditorHTML string `json:"guideEditorHtml"`

	WarrantRequestCount    string `json:"warrantRequestCount"`
	WarrantProcessingCount string `json:"warrantProcessingCount"`
	WarrantProvideCount    string `json:"warrantProvideCount"`
	WarrantAverageCount    string `json:"warrantAverageCount"`
	WarrantRate            string `json:"warrantRate"`

	CommDataRequestCount    string `json:"commDataRequestCount"`
	CommDataProcessingCount string `json:"commDataProcessingCount"`
	CommDataProvideCount    string `json:"commDataProvideCount"`
	CommDataAverageCount    string `json:"commDataAverageCount"`
	CommDataRate            string `json:"commDataRate"`

	CommRestrictionRequestCount    string `json:"commRestrictionRequestCount"`
	CommRestrictionProcessingCount string `json:"commRestrictionProcessingCount"`
	CommRestrictionProvideCount    string `json:"commRestrictionProvideCount"`
	CommRestrictionAverageCount    string `json:"commRestrictionAverageCount"`
	CommRestrictionRate            string `json:"commRestrictionRate"`

	CommConfirmationDataRequestCount    string `json:"commConfirmationDataRequestCount"`
	CommConfirmationDataProcessingCount string `json:"commConfirmationDataProcessingCount"`
	CommConfirmationDataProvideCount    string `json:"commConfirmationDataProvideCount"`
	CommConfirmationDataAverageCount    string `json:"commConfirmationDataAverageCount"`
	CommConfirmationDataRate            string `json:"commConfirmationDataRate"`
}

// fetchNaverStatistics loads the transparency-statistics page and returns its
// half-year series entries, newest first (the CMS order).
func fetchNaverStatistics(ctx context.Context, c *client.Client) ([]naverStatistic, error) {
	raw, err := c.Get(ctx, "/api/pages/"+naverStatisticsPageCode, nil)
	if err != nil {
		return nil, err
	}
	var page struct {
		SpecificAreaJSON struct {
			Statistics []naverStatistic `json:"statistics"`
		} `json:"specificAreaJson"`
	}
	if err := json.Unmarshal(raw, &page); err != nil {
		return nil, fmt.Errorf("parsing %s page payload: %w", naverStatisticsPageCode, err)
	}
	if len(page.SpecificAreaJSON.Statistics) == 0 {
		return nil, fmt.Errorf("page %s carried no statistics entries; the CMS payload shape may have changed", naverStatisticsPageCode)
	}
	return page.SpecificAreaJSON.Statistics, nil
}

// naverPeriodLabel renders the CMS's Korean half-year label as the compact
// 1H/2H form used in output rows. Unknown labels pass through untouched.
func naverPeriodLabel(period string) string {
	switch strings.TrimSpace(period) {
	case "상반기":
		return "1H"
	case "하반기":
		return "2H"
	default:
		return strings.TrimSpace(period)
	}
}

var htmlTagPattern = regexp.MustCompile(`<[^>]*>`)
var blankLinesPattern = regexp.MustCompile(`\n{3,}`)

// stripHTML flattens CMS HTML into readable plain text: block-level tags
// become newlines, remaining tags are dropped, and entities are decoded.
func stripHTML(s string) string {
	s = regexp.MustCompile(`(?i)</(p|li|ul|ol|div|h[1-6])>`).ReplaceAllString(s, "\n")
	s = regexp.MustCompile(`(?i)<br\s*/?>`).ReplaceAllString(s, "\n")
	s = htmlTagPattern.ReplaceAllString(s, "")
	s = nethtml.UnescapeString(s)
	lines := strings.Split(s, "\n")
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}
	s = strings.Join(lines, "\n")
	s = blankLinesPattern.ReplaceAllString(s, "\n\n")
	return strings.TrimSpace(s)
}
