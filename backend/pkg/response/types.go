package response

type APIError struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}