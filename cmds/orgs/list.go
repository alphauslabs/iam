package orgs

import (
	"context"
	"encoding/json"
	"io"
	"os"

	"github.com/alphauslabs/blue-internal-go/iam/v1"
	"github.com/alphauslabs/bluectl/pkg/logger"
	"github.com/alphauslabs/iam/params"
	"github.com/alphauslabs/iam/pkg/connection"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func ListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all orgs",
		Long:  `List all orgs.`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			con, err := connection.New(ctx)
			if err != nil {
				logger.Errorf("connection.New failed: %v", err)
				return
			}

			defer con.Close()
			client := iam.NewIamServiceClient(con)
			stream, err := client.ListOrgs(ctx, &iam.ListOrgsRequest{})
			if err != nil {
				logger.Errorf("ListOrgs failed: %v", err)
				return
			}

			render := true
			table := tablewriter.NewWriter(os.Stdout)
			table.SetBorder(false)
			table.SetHeaderLine(false)
			table.SetColumnSeparator("")
			table.SetTablePadding("  ")
			table.SetNoWhiteSpace(true)
			table.Append([]string{"ID", "NAME", "EMAIL", "STATUS"})

			for {
				v, err := stream.Recv()
				if err == io.EOF {
					break
				}

				switch {
				case params.OutFmt == "json":
					render = false
					b, _ := json.Marshal(v)
					logger.Info(string(b))
				default:
					var email string
					sts := "-"
					m := v.Attributes.AsMap()
					if _, ok := m["email"]; ok {
						email = m["email"].(string)
					}

					if _, ok := m["status"]; ok {
						sts = m["status"].(string)
					}

					table.Append([]string{v.Id, v.Name, email, sts})
				}
			}

			if render {
				table.Render()
			}
		},
	}

	cmd.Flags().SortFlags = false
	return cmd
}
