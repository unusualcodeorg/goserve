package micro

type sender struct{}

func NewMessageSender() MessageSender {
	return &sender{}
}

func (m *sender) SendNats(req NatsRequest) SendMessage {
	return &send{
		natsRequest: req,
	}
}

type send struct {
	natsRequest NatsRequest
}

func (s *send) Message(data any) {
	s.natsRequest.RespondJSON(NewAnyMessage(data, nil))
}

func (s *send) Error(err error) {
	s.natsRequest.RespondJSON(NewAnyMessage(nil, err))
}
