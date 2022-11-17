package users

type WrongUsernameOrPassword struct{}

func (WrongUsernameOrPassword) Error() string {
	return "wrong username or password"
}
