package cmds

import (
	"context"
	"crypto/tls"

	"github.com/alphauslabs/blue-internal-go/iam/v1"
	"github.com/alphauslabs/bluectl/pkg/logger"
	"github.com/spf13/cobra"
	"google.golang.org/api/idtoken"
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
			// dev: iamd-dev-cnugyv5cta-an.a.run.app
			// next: iamd-next-u554nqhjka-an.a.run.app
			// prod: iamd-prod-u554nqhjka-an.a.run.app
			svc := "iamd-prod-u554nqhjka-an.a.run.app"
			var opts []grpc.DialOption
			creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
			opts = append(opts, grpc.WithTransportCredentials(creds))
			opts = append(opts, grpc.WithBlock())
			ts, err := idtoken.NewTokenSource(ctx, "https://"+svc)
			if err != nil {
				logger.Error(err)
				return
			}

			token, err := ts.Token()
			if err != nil {
				logger.Error(err)
				return
			}

			opts = append(opts, grpc.WithUnaryInterceptor(func(ctx context.Context,
				method string, req, reply interface{}, cc *grpc.ClientConn,
				invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
				ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token.AccessToken)
				return invoker(ctx, method, req, reply, cc, opts...)
			}))

			opts = append(opts, grpc.WithStreamInterceptor(func(ctx context.Context,
				desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer,
				opts ...grpc.CallOption) (grpc.ClientStream, error) {
				ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token.AccessToken)
				return streamer(ctx, desc, cc, method, opts...)
			}))

			hp := svc + ":443"
			ccon, err := grpc.DialContext(ctx, hp, opts...)
			if err != nil {
				logger.Errorf("DialContext failed: %v", err)
				return
			}

			defer ccon.Close()
			client := iam.NewIamServiceClient(ccon)
			resp, err := client.WhoAmI(ctx, &iam.WhoAmIRequest{})
			if err != nil {
				logger.Error(err)
				return
			}

			logger.Info(resp)
		},
	}

	cmd.Flags().SortFlags = false
	return cmd
}
