package types

// StandardResponse Standardized API response
type StandardResponse struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}
