package users

import (
	database "github.com/masudur-rahman/hackernews/internal/pkg/db/migrations/mysql"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (user *User) Create() error {
	statement, err := database.Db.Prepare("INSERT INTO Users(Username, Password) VALUES (?,?)")
	if err != nil {
		return err
	}

	hash, err := HashPassword(user.Password)
	if err != nil {
		return err
	}

	if _, err = statement.Exec(user.Username, hash); err != nil {
		return err
	}

	return nil
}

// HashPassword hashes given password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash compares raw password with it's hashed values
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetUserIDByUsername(username string) (int, error) {
	statement, err := database.Db.Prepare("select ID from Users where Username = ?")
	if err != nil {
		return 0, err
	}

	row := statement.QueryRow(username)

	var ID int
	err = row.Scan(&ID)
	if err != nil {
		return 0, err
	}

	return ID, err
}
