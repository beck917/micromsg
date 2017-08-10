package helpers

var GlobalWorker interface{}
var GlobalUidBindClientId map[int]*UserData = make(map[int]*UserData)
var GlobalClientIdBindUid map[uint64]int = make(map[uint64]int)

type UserData struct {
	ClientId uint64
	Name     string
}

func NewUserData() *UserData {
	ud := new(UserData)
	ud.ClientId = 0

	return ud
}
