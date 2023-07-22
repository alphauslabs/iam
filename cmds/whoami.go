package cmds

import (
	"context"
	"encoding/json"

	"github.com/alphauslabs/blue-internal-go/iam/v1"
	"github.com/alphauslabs/bluectl/pkg/logger"
	"github.com/alphauslabs/iam/pkg/connection"
	"github.com/spf13/cobra"
)

func WhoAmICmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whoami",
		Short: "Get my information as an internal user",
		Long:  `Get my information as an internal user.`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			con, err := connection.New(ctx)
			if err != nil {
				logger.Errorf("connection.New failed: %v", err)
				return
			}

			defer con.Close()
			client := iam.NewIamServiceClient(con)
			resp, err := client.WhoAmI(ctx, &iam.WhoAmIRequest{})
			if err != nil {
				logger.Errorf("WhoAmI failed: %v", err)
				return
			}

			b, _ := json.Marshal(resp.Fields)
			logger.Info(string(b))
		},
	}

	cmd.Flags().SortFlags = false
	return cmd
}
