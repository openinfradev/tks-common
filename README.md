# tks-common

[![Go Report Card](https://goreportcard.com/badge/github.com/openinfradev/tks-common?style=flat-square)](https://goreportcard.com/report/github.com/openinfradev/common)
[![Go Reference](https://pkg.go.dev/badge/github.com/openinfradev/common.svg)](https://pkg.go.dev/github.com/openinfradev/common)
[![Release](https://img.shields.io/github/release/sktelecom/tks-common.svg?style=flat-square)](https://github.com/openinfradev/tks-common/releases/latest)

TKS는 TACO Kubernetes Service의 약자로, SK Telecom이 만든 GitOps기반의 서비스 시스템을 의미합니다. 그 중 tks-common은 각 tks 컴포넌트에서 공통으로 사용하는 정적인 라이브러리를 관리하는 컴포넌트입니다. 모든 tks 컴포넌트는 본 repository 를 import 하는 구조로 설계하며, 상호 의존을 지양합니다.

| pkg    | 목적         |
|--------|-------------|
| argowf | argo workflow 를 호출하는 client 를 제공 |
| grpc_client | tks-contract, tks-info, tks-cluster-lcm 을 호출하는 grpc client 를 표준화 |
| log | tks-apis 에서 사용하는 log package 로, 로그 레벨에 따른 로그출력을 표준화 |
| helper | 기타 utilities |


### tks-common을 사용한 tks service startup 예제
```go
import  "github.com/openinfradev/tks-common/pkg/grpc_server"

func YOUR_FUNCTION(YOUR_PARAMS...) {
	// start server
	s, conn, err := grpc_server.CreateServer(port, tlsEnabled, tlsCertPath, tlsKeyPath)
	if err != nil {
		log.Fatal("failed to crate grpc_server : ", err)
	}
}
```

### tks-common을 사용한 tks component 간 gRPC API 호출 예제 (golang)
```go
import pb "github.com/openinfradev/tks-proto/tks_pb"
import  "github.com/openinfradev/tks_common/pkg/gprc_client"

func YOUR_FUNCTION(YOUR_PARAMS...) {
    if _, clusterInfoClient, err = grpc_client.CreateClusterInfoClient(infoAddress, infoPort, false, ""); err != nil {
        // error
    }
    res, err := clusterInfoClient.GetCluster(ctx, &pb.GetClusterRequest{ClusterId: clusterId})
    if err != nil {
        // error
    }            
}

```

### tks-common을 사용한 tks logging 예제
```go
import  "github.com/openinfradev/tks-common/pkg/log"

func YOUR_FUNCTION(YOUR_PARAMS...) {
	log.Debug( "DEBUG 목적의 로깅을 위해 사용합니다. 서버 startup시 LOG_LEVEL을 명시적으로 DEBUG 로 설정해야 출력됩니다.")
	log.Info( "일반적인 INFO 목적의 로깅을 위해 사용합니다.")
	log.Error( "ERROR 목적의 로깅을 위해 사용합니다.")
	log.Fatal( "일반적으로 사용되지 않습니다. Fatal 로깅후 서버를 강제로 종료합니다.")

	log.Debug( "간단한 parameter 로깅은 이와 같이 표현할 수 있습니다. ", param1 )
	log.Debug( fmt.Sprintf("상세한 parameter 로깅은 fmt.Sprintf() 를 혼합하여 표현 할 수 있습니다. (%s) ", param1) )
}
```

### tks-common을 사용한 argo workflow 호출 예제
```go
import  "github.com/openinfradev/tks-common/pkg/argowf"

func YOUR_FUNCTION(YOUR_PARAMS...) {
	argowfClient, err = argowf.New(argoAddress, argoPort)
	if err != nil {
		// error
	}

	workflowId, err := argowfClient.SumbitWorkflowFromWftpl(ctx, workflow, nameSpace, parameters)
	if err != nil {
		// error
	}
}
```
