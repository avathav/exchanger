package exchange

type ApiResponse struct {
	Timestamp int64              `json:"timestamp"`
	Base      string             `json:"base"`
	Rates     map[string]float64 `json:"rates"`
}

type ResponseError struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	Description string `json:"description"`
}
