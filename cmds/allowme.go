package cmds

import (
	"context"

	"github.com/alphauslabs/blue-internal-go/iam/v1"
	"github.com/alphauslabs/bluectl/pkg/logger"
	"github.com/alphauslabs/iam/pkg/connection"
	"github.com/spf13/cobra"
)

func AllowMeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "allow-me <service-name>",
		Short: "Request access for an internal service",
		Long:  `Request access for an internal CloudRun service.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				logger.Errorf("[service-name] is required")
				return
			}

			ctx := context.Background()
			con, err := connection.New(ctx)
			if err != nil {
				logger.Errorf("connection.New failed: %v", err)
				return
			}

			defer con.Close()
			client := iam.NewIamServiceClient(con)
			resp, err := client.AllowMe(ctx, &iam.AllowMeRequest{
				Service: args[0],
			})

			if err != nil {
				logger.Errorf("AllowMe failed: %v", err)
				return
			}

			logger.Infof("%v", resp.Message)
		},
	}

	cmd.Flags().SortFlags = false
	return cmd
}
