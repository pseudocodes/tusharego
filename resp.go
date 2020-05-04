package tushare

type ProResponse struct {
	RequestID string `json:"request_id"`
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	Data      struct {
		Fields []string        `json:"fields"`
		Items  [][]interface{} `json:"items"`
	} `json:"data"`
}
