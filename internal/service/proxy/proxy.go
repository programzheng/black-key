package proxy

import (
	"context"
	"time"

	"github.com/programzheng/black-key/config"
	pb "github.com/programzheng/black-key/pkg/proto/proxy"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetGrpcProxyResponse(identifier *string, key string) string {
	grpcProxyUrl := config.Cfg.GetString("GRPC_PROXY_URL")
	if grpcProxyUrl == "" {
		return ""
	}
	conn, err := grpc.Dial(grpcProxyUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewProxyClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()
	r, err := c.GetProxy(ctx, &pb.GetProxyRequest{
		Identifier: identifier,
		Key:        key,
	})
	if err != nil {
		log.Fatalf("could not get proxy response: %v", err)
	}
	return *r.Url
}
