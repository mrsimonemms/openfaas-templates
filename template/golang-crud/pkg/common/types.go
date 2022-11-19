package common

type Paginate struct {
	Count        int `json:"count"`
	Data         any `json:"data"`
	Page         int `json:"page"`
	TotalPages   int `json:"pages"`
	TotalResults int `json:"total"`
}
