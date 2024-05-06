package entity

type ChildUser struct {
	id UserID
}

func NewChildUser(id UserID) (hu *ChildUser) {
	return hu
}

func (u *ChildUser) GetAudioNum() {

}
