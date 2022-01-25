package grpc_client

import (
	"sync"
	"fmt"

	"google.golang.org/grpc"
	//"google.golang.org/grpc/credentials"

	pb "github.com/openinfradev/tks-proto/tks_pb"
)

var (
	once sync.Once
	client pb.ClusterLcmServiceClient
)

func GetClusterLcmServiceClient(address string, port int, caller string) pb.ClusterLcmServiceClient {
	host := fmt.Sprintf("%s:%d", address, port)
	once.Do(func() {
		conn, _ := grpc.Dial(
			host,
			grpc.WithInsecure(),
		)

		client = pb.NewClusterLcmServiceClient(conn)
	})
	return client
}

