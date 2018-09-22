package pgway

import (
	"encoding/json"
	"net/http"
)

type Request struct {
	Path            string            //
	HTTPMethod      string            //
	RequestData     map[string]string //
	QueryParameters map[string]string //
	PathVariables   map[string]string //
	Headers         map[string]string //
	Body            string            //
}

//
// integrate queryData with postData
func (req *Request) initRequestData() error {

	data := make(map[string]string)

	for key, value := range req.QueryParameters {
		data[key] = value
	}

	if req.HTTPMethod == http.MethodPost {

		postData := map[string]string{}
		//
		// MARK: currentry not suported form type post data (exp a=b\n c=d,,,)
		err := json.Unmarshal([]byte(req.Body), &postData)

		if err != nil {
			return err
		}

		for key, value := range postData {
			data[key] = value
		}
	}

	for key, value := range req.PathVariables {
		data[key] = value
	}

	req.RequestData = data

	return nil
}
