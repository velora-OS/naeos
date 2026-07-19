package main

import (
	"github.com/spf13/cobra"

	"github.com/NAEOS-foundation/naeos/internal/api"
	"github.com/NAEOS-foundation/naeos/internal/dashboard"
)

var (
	dashPort string
)

func newDashboardCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dashboard",
		Short: "Start NAEOS web dashboard",
		Long:  `Start the NAEOS web dashboard for monitoring and managing projects.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			dash, err := dashboard.New()
			if err != nil {
				return err
			}

			mux := api.NewServer(dashPort, &api.AuthConfig{Enabled: false})
			mux.Router.HandleFunc("/", dash.ServeHTTP)

			return mux.Start()
		},
	}

	cmd.Flags().StringVarP(&dashPort, "port", "p", "3000", "Dashboard port")

	return cmd
}
