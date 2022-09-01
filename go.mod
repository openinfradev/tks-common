module github.com/openinfradev/tks-common

go 1.16

require (
	github.com/argoproj/argo-workflows/v3 v3.1.13
	github.com/golang/mock v1.6.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0
	github.com/lib/pq v1.10.4
	github.com/openinfradev/tks-proto v0.0.6-0.20220831015809-fad377174017
	github.com/ory/dockertest/v3 v3.8.2-0.20220202112136-e58dff82f532
	github.com/sirupsen/logrus v1.8.1
	google.golang.org/grpc v1.43.0
	k8s.io/apimachinery v0.19.6
)

replace github.com/openinfradev/tks-common => ./

//replace github.com/openinfradev/tks-proto => ../tks-proto
