package links

import (
	"log"

	database "github.com/masudur-rahman/hackernews/internal/pkg/db/migrations/mysql"
	"github.com/masudur-rahman/hackernews/internal/users"
)

type Link struct {
	ID      string
	Title   string
	Address string
	User    *users.User
}

func (link Link) Save() (int64, error) {
	stmt, err := database.Db.Prepare("INSERT INTO Links(Title, Address, UserID) VALUES (?,?,?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(link.Title, link.Address, link.User.ID)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	log.Println("Row inserted!")
	return id, nil
}

func GetAll() ([]Link, error) {
	stmt, err := database.Db.Prepare("select l.id, l.title, l.address, l.UserID, u.Username from Links l inner join Users u on l.UserID = u.ID")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []Link
	var id string
	var username string

	for rows.Next() {
		var link Link
		if err = rows.Scan(&link.ID, &link.Title, &link.Address, &id, &username); err != nil {
			return nil, err
		}

		link.User = &users.User{
			ID:       id,
			Username: username,
		}
		links = append(links, link)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return links, nil
}
