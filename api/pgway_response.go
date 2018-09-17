package api

type PgwayResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       string
}
