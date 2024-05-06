package entity

type UserMsg struct {
	MsgId string
	From  UserID
	Text  string
	Audio Audio
	Url   string
}
