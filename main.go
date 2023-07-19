package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/alphauslabs/bluectl/pkg/logger"
	"github.com/alphauslabs/iam/cmds"
	"github.com/alphauslabs/iam/params"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
)

var (
	bold = color.New(color.Bold).SprintFunc()
	year = func() string {
		return fmt.Sprintf("%v", time.Now().Year())
	}

	rootCmd = &cobra.Command{
		Use:   "iam",
		Short: bold("iam") + " - Command line interface for iamd",
		Long: bold("iam") + ` - Command line interface for our internal IAM service.
Copyright (c) 2023-` + year() + ` Alphaus Cloud, Inc. All rights reserved.

The general form is ` + bold("iam <resource[ subresource...]> <action> [flags]") + `.

To authenticate, either set GOOGLE_APPLICATION_CREDENTIALS env var or
set the --creds-file flag. Ask the service owner if your credentials
file doesn't have access to the service itself.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if params.AccessToken != "" {
				return
			}

			svc := "iamd-prod-u554nqhjka-an.a.run.app"
			switch params.RunEnv {
			case "dev":
				svc = "iamd-dev-cnugyv5cta-an.a.run.app"
			case "next":
				svc = "iamd-next-u554nqhjka-an.a.run.app"
			}

			ctx := context.Background()
			var ts oauth2.TokenSource
			var err error
			switch {
			case params.CredentialsFile != "":
				opts := idtoken.WithCredentialsFile(params.CredentialsFile)
				ts, err = idtoken.NewTokenSource(ctx, "https://"+svc, opts)
			default:
				ts, err = idtoken.NewTokenSource(ctx, "https://"+svc)
			}

			if err != nil {
				logger.Error(err)
				os.Exit(1)
			}

			token, err := ts.Token()
			if err != nil {
				logger.Error(err)
				os.Exit(1)
			}

			params.AccessToken = token.AccessToken
		},
		Run: func(cmd *cobra.Command, args []string) {
			logger.Info("see -h for more information")
		},
	}
)

func init() {
	rootCmd.Flags().SortFlags = false
	rootCmd.PersistentFlags().SortFlags = false
	rootCmd.PersistentFlags().StringVar(&params.CredentialsFile, "creds-file", "", "optional, GCP service account file")
	rootCmd.PersistentFlags().StringVar(&params.AccessToken, "access-token", "", "use directly if not empty")
	rootCmd.PersistentFlags().StringVar(&params.RunEnv, "env", "prod", "dev, next, or prod")
	rootCmd.AddCommand(
		cmds.WhoAmICmd(),
	)
}

func main() {
	cobra.EnableCommandSorting = false
	log.SetOutput(os.Stdout)
	rootCmd.Execute()
}
