package api

type PgwayApis []PgwayApi

func (apis PgwayApis) Len() int {
	return len(apis)
}

//
// 最長一致
func (apis PgwayApis) Less(i, j int) bool {
	return len(apis[i].Path) > len(apis[j].Path)
}

// Swap swaps the elements with indexes i and j.package github.com/awslabs/aws-lambda-go-api-proxy/...: cannot download, $GOPATH must not be set to $GOROOT. For more details see: 'go help gopath'
func (apis PgwayApis) Swap(i, j int) {
	apis[i], apis[j] = apis[j], apis[i]
}
