package pgway

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type PgwayHttpConf struct {
	Port int
}

func (req *Request) WitHttpRequest(baseRequest *http.Request) {
	req.Path = baseRequest.URL.Path
	req.HTTPMethod = baseRequest.Method

	//
	// query parameters
	req.QueryParameters = map[string]string{}
	for key, value := range baseRequest.URL.Query() {
		req.QueryParameters[key] = strings.Join(value, ",")
	}
	//
	// header
	req.Headers = map[string]string{}
	for key, value := range baseRequest.Header {
		//
		// why value is array?
		req.Headers[key] = strings.Join(value, ",")
	}

	//
	// body
	b, err := ioutil.ReadAll(baseRequest.Body)
	if err != nil {
		panic(err)
	}
	req.Body = string(b)
}

func (resp *Response) WriteHttpResponse(w http.ResponseWriter) {

	//
	// status code
	w.WriteHeader(resp.StatusCode)
	//
	// header
	for key, value := range resp.Headers {
		w.Header().Add(key, value)
	}

	//
	// body
	w.Write([]byte(resp.Body))

}

func (server Server) ServeHTTP(w http.ResponseWriter, httpRequest *http.Request) {
	request := &Request{}
	request.WitHttpRequest(httpRequest)
	response := server.handle(request)
	response.WriteHttpResponse(w)
}

func (server *Server) BootHttpServerWithDefaultConfig() {
	conf := PgwayHttpConf{
		Port: 8080,
	}
	server.BootHttpServer(conf)
}

func (server *Server) BootHttpServer(conf PgwayHttpConf) {

	addr := ":" + strconv.Itoa(conf.Port)
	http.ListenAndServe(addr, server)
}
