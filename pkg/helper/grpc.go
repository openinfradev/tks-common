package helper

import (
	"fmt"

	"google.golang.org/grpc"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"

	"github.com/openinfradev/tks-common/pkg/log"
)

func CreateConnection(address string, port int, caller string) (*grpc.ClientConn, error) {
	host := fmt.Sprintf("%s:%d", address, port)
	conn, err := grpc.Dial(
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
	return conn, nil
}

