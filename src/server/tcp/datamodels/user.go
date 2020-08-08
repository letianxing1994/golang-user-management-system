package datamodels

type User struct {
	UserId         int64  `sql:"user_id"`
	Nickname       string `sql:"nickname"`
	ProfilePicture string `sql:"profile_picture"`
	Password       string `sql:"password"`
	Username       string `sql:"username"`
}
