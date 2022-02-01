package sms

type Item struct {
	Number    string `json:"number"`
	Text      string `json:"message"`
	RequestId string `json:"request_id"`
}

func (i Item) Phone() string {
	return i.Number
}
func (i Item) Message() string {
	return i.Text
}
func (i Item) ID() string {
	return i.RequestId
}
