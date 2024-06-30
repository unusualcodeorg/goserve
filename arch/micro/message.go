package micro

type Message[T any] struct {
	Data  *T      `json:"data,omitempty"`
	Error *string `json:"error,omitempty"`
}

func NewMessage[T any](data *T, err error) *Message[T] {
	var e *string
	if err != nil {
		er := err.Error()
		e = &er
	}

	return &Message[T]{
		Data:  data,
		Error: e,
	}
}
