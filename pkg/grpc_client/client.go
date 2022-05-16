package grpc_client

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/openinfradev/tks-common/pkg/log"
	pb "github.com/openinfradev/tks-proto/tks_pb"
)

func CreateCspInfoClient(address string, port int, tlsEnabled bool, certPath string) (*grpc.ClientConn, pb.CspInfoServiceClient, error) {
	cc, err := createConnection(address, port, tlsEnabled, certPath)
	if err != nil {
		log.Fatal("Could not connect to gRPC server", err)
		return nil, nil, err
	}
	sc := pb.NewCspInfoServiceClient(cc)
	return cc, sc, nil
}

func CreateContractClient(address string, port int, tlsEnabled bool, certPath string) (*grpc.ClientConn, pb.ContractServiceClient, error) {
	cc, err := createConnection(address, port, tlsEnabled, certPath)
	if err != nil {
		log.Fatal("Could not connect to gRPC server", err)
		return nil, nil, err
	}
	sc := pb.NewContractServiceClient(cc)
	return cc, sc, nil
}

func CreateClusterInfoClient(address string, port int, tlsEnabled bool, certPath string) (*grpc.ClientConn, pb.ClusterInfoServiceClient, error) {
	cc, err := createConnection(address, port, tlsEnabled, certPath)
	if err != nil {
		log.Fatal("Could not connect to gRPC server", err)
		return nil, nil, err
	}
	sc := pb.NewClusterInfoServiceClient(cc)
	return cc, sc, nil
}

func CreateAppInfoClient(address string, port int, tlsEnabled bool, certPath string) (*grpc.ClientConn, pb.AppInfoServiceClient, error) {
	cc, err := createConnection(address, port, tlsEnabled, certPath)
	if err != nil {
		log.Fatal("Could not connect to gRPC server", err)
		return nil, nil, err
	}
	sc := pb.NewAppInfoServiceClient(cc)
	return cc, sc, nil
}

func createConnection(address string, port int, tlsEnabled bool, certPath string) (*grpc.ClientConn, error) {
	var err error
	var creds credentials.TransportCredentials

	if tlsEnabled {
		creds, err = loadTLSClientCredential(certPath)
		if err != nil {
			return nil, err
		}
	} else {
		creds = insecure.NewCredentials()
	}

	host := fmt.Sprintf("%s:%d", address, port)
	conn, err := grpc.Dial(
		host,
		grpc.WithTransportCredentials(creds),
		grpc.WithUnaryInterceptor(
			grpc_middleware.ChainUnaryClient(
				log.IOLoggingForClientSide(),
			),
		),
	)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func loadTLSClientCredential(clientCertPath string) (credentials.TransportCredentials, error) {
	creds, err := credentials.NewClientTLSFromFile(clientCertPath, "")
	if err != nil {
		log.Error("Fail to load client credentials: ", err)
		return nil, err
	}

	return creds, nil
}
