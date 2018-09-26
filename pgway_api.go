package pgway

type Api struct {
	Path       string      // api path
	HTTPMethod string      // httpMethod
	Handler    interface{} // api method handler
	//IsDebug    bool        // debug api // TODO: develop
}
