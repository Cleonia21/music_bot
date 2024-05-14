package entity

const (
	UserMsgExit = 1
)

type UserMsg struct {
	MsgId int
	From  UserID
	Text  string
	Audio Audio
	Url   string
}
