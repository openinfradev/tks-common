package log

import (
	"os"
	"strings"
	"io/ioutil"
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	logger *logrus.Logger
)

// Init initializes logrus.logger and set
func init() {
	logger = logrus.New()

	logger.Out = os.Stdout

	formatter := new(logrus.TextFormatter)
	formatter.FullTimestamp = true
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	logger.SetFormatter(formatter)

	logLevel := strings.ToLower(os.Getenv("LOG_LEVEL"))
	switch logLevel {
		case "debug":
			logger.SetLevel(logrus.DebugLevel)
		case "warning":
			logger.SetLevel(logrus.WarnLevel)
		case "info":
			logger.SetLevel(logrus.InfoLevel)
		case "error":
			logger.SetLevel(logrus.ErrorLevel)
		case "fatal":
			logger.SetLevel(logrus.FatalLevel)
		default:
			logger.SetLevel(logrus.InfoLevel)
	}
}

// Info logs in InfoLevel.
func Info(v ...interface{}) {
	logger.Info(v...)
}

// Warn logs in WarnLevel.
func Warn(v ...interface{}) {
	logger.Warn(v...)
}

// Debug logs in DebugLevel.
func Debug(v ...interface{}) {
	logger.Debug(v...)
}

// Error logs in ErrorLevel.
func Error(v ...interface{}) {
	logger.Error(v...)
}

// Fatal logs in FatalLevel
func Fatal(v ...interface{}) {
	logger.Fatal(v...)
}

func Disable() {
	logger.Out = ioutil.Discard
}


// grpc IO logging for client-side 
func IOLoggingForClientSide() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		err := invoker(ctx, method, req, reply, cc, opts...)

		Info(fmt.Sprintf("[INTERNAL_CALL:%s][REQUEST %s][RESPONSE %s]", method, req, reply))
		
		return err
	}
}

// grpc IO logging for server-side 
func IOLoggingForServerSide() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		Info(fmt.Sprintf("[START:%s][REQUEST %s]", info.FullMethod, req))

		res, err := handler(ctx, req)
		if err != nil {
			Error(err)
		}

		Info(fmt.Sprintf("[END:%s][RESPONSE %s]", info.FullMethod, res))

		return res, err
	}
}
