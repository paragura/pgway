package api

type PgwayApi struct {
	Path       string      // api path
	HTTPMethod string      // httpMethod
	Handler    interface{} // api method handler
	//IsDebug    bool        // debug api // TODO: develop
}

//
// TODO: path variable
func (api PgwayApi) IsSamePath(path string) bool {
	return path == api.Path
}
