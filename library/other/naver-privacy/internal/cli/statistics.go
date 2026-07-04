// Copyright 2026 Kieran Maynard and contributors. Licensed under Apache-2.0. See LICENSE.
// Hand-authored novel feature for the Naver Privacy CLI. Carried across regen
// via the novel-command merge path; the generated stub was replaced.

package cli

import (
	"strings"

	"github.com/spf13/cobra"
)

// statisticsRow is one half-year x request-category observation reshaped from
// the CMS's 24-wide-column entries. Counts stay strings for source fidelity:
// "-" marks categories Naver no longer complies with (no communication user
// information provided since October 2012).
type statisticsRow struct {
	Year       int    `json:"year"`
	Period     string `json:"period"`
	Category   string `json:"category"`
	CategoryKo string `json:"categoryKo"`
	Requests   string `json:"requests"`
	Processed  string `json:"processed"`
	Provided   string `json:"provided"`
	Average    string `json:"average"`
	Rate       string `json:"rate"`
}

// naverCategories decodes the CMS's per-category column prefixes into stable
// category slugs and the official Korean names.
var naverCategories = []struct {
	slug string
	ko   string
	pick func(s naverStatistic) (requests, processed, provided, average, rate string)
}{
	{"warrant", "압수수색영장", func(s naverStatistic) (string, string, string, string, string) {
		return s.WarrantRequestCount, s.WarrantProcessingCount, s.WarrantProvideCount, s.WarrantAverageCount, s.WarrantRate
	}},
	{"comm-user-info", "통신이용자정보", func(s naverStatistic) (string, string, string, string, string) {
		return s.CommDataRequestCount, s.CommDataProcessingCount, s.CommDataProvideCount, s.CommDataAverageCount, s.CommDataRate
	}},
	{"comm-restriction", "통신제한조치", func(s naverStatistic) (string, string, string, string, string) {
		return s.CommRestrictionRequestCount, s.CommRestrictionProcessingCount, s.CommRestrictionProvideCount, s.CommRestrictionAverageCount, s.CommRestrictionRate
	}},
	{"comm-confirmation", "통신사실확인자료", func(s naverStatistic) (string, string, string, string, string) {
		return s.CommConfirmationDataRequestCount, s.CommConfirmationDataProcessingCount, s.CommConfirmationDataProvideCount, s.CommConfirmationDataAverageCount, s.CommConfirmationDataRate
	}},
}

// pp:data-source live
func newNovelStatisticsCmd(flags *rootFlags) *cobra.Command {
	var flagCategory string
	var flagYear int

	cmd := &cobra.Command{
		Use:   "statistics",
		Short: "The full 2012-present transparency series as tidy rows: one per half-year and request category",
		Long: "Fetches the transparency-statistics page and reshapes its wide half-year entries into " +
			"long-form rows — one per half-year and government request category (search & seizure warrants, " +
			"communication user information, communication-restricting measures, communication confirmation " +
			"data) — with request/processed/provided/average counts and compliance rates. Counts are " +
			"source-fidelity strings; \"-\" means not applicable.",
		Example: "  naver-privacy-pp-cli statistics --category warrant\n" +
			"  naver-privacy-pp-cli statistics --year 2025 --agent --select year,period,category,requests,processed,rate\n" +
			"  naver-privacy-pp-cli statistics --csv",
		Annotations: map[string]string{"mcp:read-only": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if dryRunOK(flags) {
				return nil
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

			wantCategory := strings.ToLower(strings.TrimSpace(flagCategory))
			rows := make([]statisticsRow, 0, len(stats)*len(naverCategories))
			for _, s := range stats {
				if flagYear != 0 && s.Year != flagYear {
					continue
				}
				for _, cat := range naverCategories {
					if wantCategory != "" && !strings.Contains(cat.slug, wantCategory) {
						continue
					}
					requests, processed, provided, average, rate := cat.pick(s)
					rows = append(rows, statisticsRow{
						Year:       s.Year,
						Period:     naverPeriodLabel(s.Period),
						Category:   cat.slug,
						CategoryKo: cat.ko,
						Requests:   requests,
						Processed:  processed,
						Provided:   provided,
						Average:    average,
						Rate:       rate,
					})
				}
			}
			return printJSONFiltered(cmd.OutOrStdout(), rows, flags)
		},
	}
	cmd.Flags().StringVar(&flagCategory, "category", "", "Filter to one request category by slug substring: warrant, comm-user-info, comm-restriction, comm-confirmation")
	cmd.Flags().IntVar(&flagYear, "year", 0, "Filter to one report year (0 = all years)")
	return cmd
}
