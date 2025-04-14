package database

type DatabaseInterface interface {
	LoginUser(username, password string) string
	LogoutUser(username, password string) bool
	RegisterUser(username, password string) string
	SetupDatabase() error
}

/*func NewDatabase() (*DatabaseInterface, error) {
	var db DatabaseInterface =
}*/
