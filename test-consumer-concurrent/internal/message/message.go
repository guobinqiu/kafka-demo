package message

type Message struct {
	ID      string `json:"id"`
	OrderID string `json:"orderId"`
	Content string `json:"content"`
}
