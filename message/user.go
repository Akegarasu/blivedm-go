package message

type User struct {
	Uid          int
	Uname        string
	Admin        bool
	Urank        int
	MobileVerify bool
	Medal        *Medal
	GuardLevel   int
	UserLevel    int64
}

type Medal struct {
	Name     string
	Level    int
	Color    int
	UpRoomId int
	UpUid    int
	UpName   string
}
