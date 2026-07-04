// Copyright 2026 Kieran Maynard and contributors. Licensed under Apache-2.0. See LICENSE.
// Hand-authored novel feature for the Naver Privacy CLI. Carried across regen
// via the novel-command merge path; the generated stub was replaced.

package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

// whitepaperRow is one publication in the NAVER Privacy Report / whitepaper
// library, flattened from the CMS's nested groups.
type whitepaperRow struct {
	Group        string `json:"group"`
	Title        string `json:"title"`
	Author       string `json:"author,omitempty"`
	DocumentPath string `json:"documentPath,omitempty"`
	DocumentName string `json:"documentName,omitempty"`
}

// pp:data-source live
func newNovelWhitepapersCmd(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whitepapers",
		Short: "Index the NAVER Privacy Report and privacy-whitepaper publications with their attached documents",
		Long: "The publication library lives inside the NAVER_REPORT page's specificAreaJson as nested " +
			"groups. This command flattens it into one searchable list of titles, authors, and attached " +
			"document paths — the primary sources behind the statistics.",
		Example:     "  naver-privacy-pp-cli whitepapers --agent\n  naver-privacy-pp-cli whitepapers --csv",
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
			raw, err := c.Get(ctx, "/api/pages/"+naverReportPageCode, nil)
			if err != nil {
				return classifyAPIError(err, flags)
			}

			var page struct {
				SpecificAreaJSON map[string][]struct {
					Type     string `json:"type"`
					Contents []struct {
						Title    string `json:"title"`
						Author   string `json:"author"`
						Document *struct {
							Path         string `json:"path"`
							OriginalName string `json:"originalName"`
						} `json:"document"`
					} `json:"contents"`
				} `json:"specificAreaJson"`
			}
			if err := json.Unmarshal(raw, &page); err != nil {
				return apiErr(fmt.Errorf("parsing %s page payload: %w", naverReportPageCode, err))
			}

			var rows []whitepaperRow
			for group, entries := range page.SpecificAreaJSON {
				for _, entry := range entries {
					for _, content := range entry.Contents {
						row := whitepaperRow{
							Group:  group,
							Title:  content.Title,
							Author: content.Author,
						}
						if content.Document != nil {
							row.DocumentPath = content.Document.Path
							row.DocumentName = content.Document.OriginalName
						}
						rows = append(rows, row)
					}
				}
			}
			if len(rows) == 0 {
				return apiErr(fmt.Errorf("page %s carried no publications; the CMS payload shape may have changed", naverReportPageCode))
			}
			return printJSONFiltered(cmd.OutOrStdout(), rows, flags)
		},
	}
	return cmd
}
