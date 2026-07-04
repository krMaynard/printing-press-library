// Copyright 2026 Kieran Maynard and contributors. Licensed under Apache-2.0. See LICENSE.
// Hand-authored novel feature for the Naver Privacy CLI. Carried across regen
// via the novel-command merge path; the generated stub was replaced.

package cli

import (
	"strconv"

	"github.com/spf13/cobra"
)

// summaryRow is one half-year's narrative commentary, HTML-stripped.
type summaryRow struct {
	Year    int    `json:"year"`
	Period  string `json:"period"`
	Summary string `json:"summary"`
}

// pp:data-source live
func newNovelSummariesCmd(flags *rootFlags) *cobra.Command {
	var flagYear string

	cmd := &cobra.Command{
		Use:   "summaries",
		Short: "Naver's own per-half-year commentary on the statistics, HTML-stripped to plain text",
		Long: "Each statistics entry carries a guideEditorHtml narrative in which Naver explains the " +
			"period's anomalies (mergers, warrant spikes, policy changes). This command extracts that " +
			"commentary as readable plain text next to the period it explains — check it before " +
			"interpreting a spike or drop in the numbers.",
		Example: "  naver-privacy-pp-cli summaries --year 2025\n" +
			"  naver-privacy-pp-cli summaries --agent --select year,period",
		Annotations: map[string]string{"mcp:read-only": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if dryRunOK(flags) {
				return nil
			}
			var wantYear int
			if flagYear != "" {
				y, err := strconv.Atoi(flagYear)
				if err != nil {
					return usageErr(err)
				}
				wantYear = y
			}
			ctx, cancel := boundCtx(cmd.Context(), flags)
			defer cancel()
			c, err := flags.newClient()
			if err != nil {
				return err
			}
			stats, err := fetchNaverStatistics(ctx, c)
			if err != nil {
				return classifyAPIError(err, flags)
			}

			rows := make([]summaryRow, 0, len(stats))
			for _, s := range stats {
				if wantYear != 0 && s.Year != wantYear {
					continue
				}
				summary := stripHTML(s.GuideEditorHTML)
				if summary == "" {
					continue
				}
				rows = append(rows, summaryRow{
					Year:    s.Year,
					Period:  naverPeriodLabel(s.Period),
					Summary: summary,
				})
			}
			return printJSONFiltered(cmd.OutOrStdout(), rows, flags)
		},
	}
	cmd.Flags().StringVar(&flagYear, "year", "", "Filter to one report year (default: all years)")
	return cmd
}
