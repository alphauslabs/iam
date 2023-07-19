package cmds

import (
	"context"
	"crypto/tls"

	"github.com/alphauslabs/blue-internal-go/iam/v1"
	"github.com/alphauslabs/bluectl/pkg/logger"
	"github.com/alphauslabs/iam/params"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

func WhoAmICmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whoami",
		Short: "Get my information as a user",
		Long:  `Get my information as a user.`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			var opts []grpc.DialOption
			creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
			opts = append(opts, grpc.WithTransportCredentials(creds))
			opts = append(opts, grpc.WithBlock())
			opts = append(opts, grpc.WithUnaryInterceptor(func(ctx context.Context,
				method string, req, reply interface{}, cc *grpc.ClientConn,
				invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
				ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+params.AccessToken)
				return invoker(ctx, method, req, reply, cc, opts...)
			}))

			opts = append(opts, grpc.WithStreamInterceptor(func(ctx context.Context,
				desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer,
				opts ...grpc.CallOption) (grpc.ClientStream, error) {
				ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+params.AccessToken)
				return streamer(ctx, desc, cc, method, opts...)
			}))

			con, err := grpc.DialContext(ctx, params.ServiceHost+":443", opts...)
			if err != nil {
				logger.Errorf("DialContext failed: %v", err)
				return
			}

			defer con.Close()
			client := iam.NewIamServiceClient(con)
			resp, err := client.WhoAmI(ctx, &iam.WhoAmIRequest{})
			if err != nil {
				logger.Error(err)
				return
			}

			for k, v := range resp.Info {
				logger.Infof("%v: %v", k, v)
			}
		},
	}

	cmd.Flags().SortFlags = false
	return cmd
}
