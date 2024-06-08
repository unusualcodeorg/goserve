package dto

type CreateMessage struct {
	Type string `json:"type" binding:"required,min=2,max=50"`
	Msg  string `json:"msg" binding:"required,min=0,max=2000"`
}
