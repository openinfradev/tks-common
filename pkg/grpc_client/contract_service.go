package grpc_client

import (
	"fmt"

	"google.golang.org/grpc"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"

	"github.com/openinfradev/tks-common/pkg/log"
	pb "github.com/openinfradev/tks-proto/tks_pb"
)

var (
	conn *grpc.ClientConn
	contractClient pb.ContractServiceClient
)

func getConnection(host string) (*grpc.ClientConn, error) {
	if conn == nil {
		_conn, err := grpc.Dial(
			host,
			grpc.WithInsecure(),
			grpc.WithUnaryInterceptor(
				grpc_middleware.ChainUnaryClient(
					log.IOLog(),
				),
			),
		)
		if err != nil {
			return nil, err
		}
		conn = _conn
	} 
	return conn, nil
}


func GetContractClient(address string, port int, caller string) (pb.ContractServiceClient, error) {
	conn, err := getConnection( fmt.Sprintf("%s:%d", address, port) )
	if err != nil {
		return nil, err
	} 

	contractClient = pb.NewContractServiceClient(conn)
	return contractClient, nil
}

