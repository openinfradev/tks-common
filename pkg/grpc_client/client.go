package grpc_client

import (
	"google.golang.org/grpc"

	"github.com/openinfradev/tks-common/pkg/log"
	"github.com/openinfradev/tks-common/pkg/helper"
	pb "github.com/openinfradev/tks-proto/tks_pb"
)

func CreateCspInfoClient(address string, port int, caller string) (*grpc.ClientConn, pb.CspInfoServiceClient, error) {
	cc, err := helper.CreateConnection(address, port, caller)
	if err != nil {
		log.Fatal("Could not connect to gRPC server", err)
		return nil, nil, err
	}
	sc := pb.NewCspInfoServiceClient(cc)
	return cc, sc, nil
}

func CreateContractClient(address string, port int, caller string) (*grpc.ClientConn, pb.ContractServiceClient, error) {
	cc, err := helper.CreateConnection(address, port, caller)
	if err != nil {
		log.Fatal("Could not connect to gRPC server", err)
		return nil, nil, err
	}
	sc := pb.NewContractServiceClient(cc)
	return cc, sc, nil
}

func CreateClusterInfoClient(address string, port int, caller string) (*grpc.ClientConn, pb.ClusterInfoServiceClient, error) {
	cc, err := helper.CreateConnection(address, port, caller)
	if err != nil {
		log.Fatal("Could not connect to gRPC server", err)
		return nil, nil, err
	}
	sc := pb.NewClusterInfoServiceClient(cc)
	return cc, sc, nil
}

func CreateAppInfoClient(address string, port int, caller string) (*grpc.ClientConn, pb.AppInfoServiceClient, error) {
	cc, err := helper.CreateConnection(address, port, caller)
	if err != nil {
		log.Fatal("Could not connect to gRPC server", err)
		return nil, nil, err
	}
	sc := pb.NewAppInfoServiceClient(cc)
	return cc, sc, nil
}

