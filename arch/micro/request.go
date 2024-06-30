package micro

import (
	"encoding/json"
	"errors"
	"fmt"
)

func Respond[T any](req NatsRequest, data *T, err error) {
	req.RespondJSON(NewMessage(data, err))
}

func Request[T any, V any](ctx *NatsContext, sub string, send *T, receive *V) (*V, error) {
	sendMsg := NewMessage(send, nil)
	sendPayload, err := json.Marshal(sendMsg)
	if err != nil {
		return nil, err
	}

	subject := fmt.Sprintf(`%s.%s`, ctx.Subject, sub)
	msg, err := ctx.Client.Conn.Request(subject, sendPayload, ctx.Client.Timeout)
	if err != nil {
		return nil, err
	}

	var receiveMsg Message[V]
	err = json.Unmarshal(msg.Data, &receiveMsg)
	if err != nil {
		return nil, err
	}

	if receiveMsg.Error != nil {
		err = errors.New(*receiveMsg.Error)
	}

	return receiveMsg.Data, err
}
