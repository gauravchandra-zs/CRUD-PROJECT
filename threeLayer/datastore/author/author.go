package datastoreauthor

import (
	"database/sql"
	"errors"

	"Projects/GoLang-Interns-2022/threeLayer/driver"
	"Projects/GoLang-Interns-2022/threeLayer/models"
)

type AuthorStore struct {
	db *sql.DB
}

func (a AuthorStore) PostAuthor(author models.Author) (int, error) {
	if !checkAuthor(author, a.db) {
		result, err := a.db.Exec(driver.InsertIntoAuthor, author.FirstName, author.LastName, author.Dob, author.PenName)
		if err != nil {
			return 0, err
		}

		id, err := result.LastInsertId()
		if err != nil {
			return 0, err
		}

		return int(id), nil
	}

	var id int

	res, err := a.db.Query(driver.CheckAuthor, author.FirstName, author.LastName, author.Dob, author.PenName)
	if err != nil {
		return 0, err
	}

	if res.Next() {
		err := res.Scan(&id)
		if err != nil {
			return 0, err
		}
	}

	return id, errors.New("author present already")
}

func (a AuthorStore) PutAuthor(id int, author models.Author) (models.Author, error) {
	res, err := a.db.Query(driver.CheckAuthorBYID, id)
	if !res.Next() || err != nil {
		return models.Author{}, errors.New("author not present")
	}

	_, err = a.db.Exec(driver.UpdateAuthor, author.FirstName, author.LastName, author.Dob, author.PenName, id)
	if err != nil {
		return models.Author{}, err
	}

	return author, nil
}

func New(db *sql.DB) AuthorStore {
	return AuthorStore{
		db,
	}
}

func (a AuthorStore) DeleteAuthor(id int) (int, error) {
	if !checkAuthorByID(id, a.db) {
		return 0, errors.New("author not exist")
	}

	res, err := a.db.Exec(driver.DeleteAuthorQuery, id)
	if err != nil {
		return 0, err
	}

	deleteID, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(deleteID), nil
}

func (a AuthorStore) GetAuthorByID(id int) (models.Author, error) {
	ResAuthor, err := a.db.Query(driver.SelectAuthorByID, id)
	if err != nil {
		return models.Author{}, err
	}

	author := models.Author{}

	if ResAuthor.Next() {
		err = ResAuthor.Scan(&author.ID, &author.FirstName, &author.LastName, &author.Dob, &author.PenName)
		if err != nil {
			return models.Author{}, err
		}
	}

	return author, nil
}

func checkAuthor(author models.Author, db *sql.DB) bool {
	res, err := db.Query(driver.CheckAuthor, author.FirstName, author.LastName, author.Dob, author.PenName)
	if !res.Next() || err != nil {
		return false
	}

	return true
}

func checkAuthorByID(id int, db *sql.DB) bool {
	res, err := db.Query(driver.CheckAuthorBYID, id)
	if !res.Next() || err != nil {
		return false
	}

	return true
}
