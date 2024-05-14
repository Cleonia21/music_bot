package entity

type UserID struct {
	ID int
}

func NewUserId(num int) UserID {
	return UserID{ID: num}
}
