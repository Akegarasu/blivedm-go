package message

type User struct {
	Uid   int
	Uname string
	Medal *Medal
}

type Medal struct {
	Name     string
	Level    int
	UpRoomId int
	UpUid    int
	UpName   string
}
