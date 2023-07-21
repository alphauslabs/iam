package cmds

import (
	"context"

	"github.com/alphauslabs/blue-internal-go/iam/v1"
	"github.com/alphauslabs/bluectl/pkg/logger"
	"github.com/alphauslabs/iam/pkg/connection"
	"github.com/spf13/cobra"
)

func WhoIsCmd() *cobra.Command {
	var (
		listAccts  bool
		listPayers bool
		listComps  bool
	)

	cmd := &cobra.Command{
		Use:   "whois <prefix:id>",
		Short: "Get information of an entity",
		Long: `Get information of an entity.

Supported prefixes:
  acct|link: - linked account
  org|msp:   - MSP
  comp:      - company
  bg:        - billing group
  rroot:     - Ripple/Octo rootuser
  rsub:      - Ripple/Octo subuser
  root:      - Wave[Pro] rootuser
  sub:       - Wave[Pro] subuser`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				logger.Errorf("<prefix:id> cannot be empty")
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
			req := iam.WhoIsRequest{
				Id:            args[0],
				ListAccts:     listAccts,
				ListPayers:    listPayers,
				ListCompanies: listComps,
			}

			resp, err := client.WhoIs(ctx, &req)
			if err != nil {
				logger.Errorf("WhoAmI failed: %v", err)
				return
			}

			b, err := resp.MarshalJSON()
			if err != nil {
				logger.Error(err)
				return
			}

			logger.Info(string(b))
		},
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().BoolVar(&listAccts, "list-accts", listAccts, "include list of accounts")
	cmd.Flags().BoolVar(&listPayers, "list-payers", listPayers, "include list of payers")
	cmd.Flags().BoolVar(&listComps, "list-companies", listComps, "include list of companies")
	return cmd
}
