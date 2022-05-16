package grpc_server

import (
	"net"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"

	"github.com/openinfradev/tks-common/pkg/log"
)

func CreateServer(port int, tlsEnabled bool, certPath string, keyPath string) (*grpc.Server, net.Listener, error) {
	log.Info("Starting to listen port ", port)

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Error("failed to listen:", err)
		return nil, nil, err
	}

	serverOptions := []grpc.ServerOption{
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_recovery.UnaryServerInterceptor(),
				log.IOLoggingForServerSide(),
			),
		),
	}

	if tlsEnabled {
		log.Info("TLS enabled!!!")
		tlsCredentials, err := loadTLSCredentials(certPath, keyPath)
		if err != nil {
			log.Error("Cannot load TLS credentials: ", err)
			return nil, nil, err
		}
		serverOptions = append(serverOptions, grpc.Creds(tlsCredentials))
	}

	return grpc.NewServer(serverOptions...), lis, nil
}

func loadTLSCredentials(certPath string, keyPath string) (credentials.TransportCredentials, error) {
	creds, err := credentials.NewServerTLSFromFile(certPath, keyPath)
	if err != nil {
		log.Error("Fail to load credentials: ", err)
		return nil, err
	}

	return creds, nil
}
