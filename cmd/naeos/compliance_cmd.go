package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/NAEOS-foundation/naeos/internal/audit"
)

func newComplianceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compliance",
		Short: "Compliance reporting and audit log export",
	}

	cmd.AddCommand(newComplianceExportCommand())
	return cmd
}

func newComplianceExportCommand() *cobra.Command {
	var format, output string

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export audit log for compliance reporting",
		Long: `Export the audit trail in JSON or CSV format for compliance purposes.

Example:
  naeos compliance export --format json --output audit-export.json
  naeos compliance export --format csv --output audit-export.csv`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if format != "json" && format != "csv" {
				return fmt.Errorf("unsupported format %q, use json or csv", format)
			}
			if output == "" {
				return fmt.Errorf("--output is required")
			}

			auditor := audit.NewMemoryAuditor()

			switch format {
			case "json":
				if err := auditor.ExportJSON(output); err != nil {
					return fmt.Errorf("export json: %w", err)
				}
			case "csv":
				if err := auditor.ExportCSV(output); err != nil {
					return fmt.Errorf("export csv: %w", err)
				}
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Compliance report exported to %s (format: %s)\n", output, format)
			return nil
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "json", "export format: json or csv")
	cmd.Flags().StringVarP(&output, "output", "o", "", "output file path (required)")
	_ = cmd.MarkFlagRequired("output")
	return cmd
}


