package dto

type CreateMessage struct {
	Type string `json:"type" binding:"required"`
	Msg  string `json:"msg" binding:"required"`
}
