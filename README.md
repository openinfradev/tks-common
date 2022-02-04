# tks-common

tks 서비스에서 사용하는 정적인 라이브러리들을 관리하는 repository.
각 tks 서비스는 상호 의존을 최대한 지양하며, tks-common 을 import 하는 구조로 설계한다.

| pkg    | 목적         |
|--------|-------------|
| argowf | argo workflow 를 호출하는 client 를 제공 |
| grpc_client | tks-contract, tks-info, tks-cluster-lcm 을 호출하는 grpc client 를 표준화 |
| log | tks-apis 에서 사용하는 log package 로, 로그 레벨에 따른 로그출력을 표준화 |
| helper | 기타 utilities |

