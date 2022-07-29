package datastoreauthor

import (
	"database/sql"

	"developer.zopsmart.com/go/gofr/pkg/gofr"

	"Projects/GoLang-Interns-2022/threeLayer/drivers"
	"Projects/GoLang-Interns-2022/threeLayer/models"
)

type AuthorStore struct {
	db *sql.DB
}

func New(db *sql.DB) AuthorStore {
	return AuthorStore{
		db,
	}
}

// PostAuthor create a row new row in author table with given detail if row not exists
func (a AuthorStore) PostAuthor(ctx *gofr.Context, author models.Author) (int, error) {
	result, err := a.db.ExecContext(ctx, drivers.InsertIntoAuthor, author.FirstName, author.LastName, author.Dob, author.PenName)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// PutAuthor update the detail of author row with given id
func (a AuthorStore) PutAuthor(ctx *gofr.Context, id int, author models.Author) (models.Author, error) {
	_, err := a.db.ExecContext(ctx, drivers.UpdateAuthor, author.FirstName, author.LastName, author.Dob, author.PenName, id)
	if err != nil {
		return models.Author{}, err
	}

	return author, nil
}

// DeleteAuthor delete author detail associated with given id if exists
func (a AuthorStore) DeleteAuthor(ctx *gofr.Context, id int) (int, error) {
	res, err := a.db.ExecContext(ctx, drivers.DeleteAuthorQuery, id)
	if err != nil {
		return 0, err
	}

	deleteID, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(deleteID), nil
}

// GetAuthorByID fetch and return author detail with given id
func (a AuthorStore) GetAuthorByID(ctx *gofr.Context, id int) (models.Author, error) {
	resAuthor, err := a.db.QueryContext(ctx, drivers.SelectAuthorByID, id)
	if err != nil {
		return models.Author{}, err
	}

	author := models.Author{}

	if resAuthor.Next() {
		err = resAuthor.Scan(&author.ID, &author.FirstName, &author.LastName, &author.Dob, &author.PenName)
		if err != nil {
			return models.Author{}, err
		}
	}

	return author, nil
}

// CheckAuthor  check author exist or not with author detail  and return bool value
func (a AuthorStore) CheckAuthor(ctx *gofr.Context, author models.Author) bool {
	row, err := a.db.QueryContext(ctx, drivers.CheckAuthor, author.FirstName, author.LastName, author.Dob, author.PenName)
	if err != nil || !row.Next() {
		return false
	}

	return true
}

// CheckAuthorByID check author exist or not with given id
func (a AuthorStore) CheckAuthorByID(ctx *gofr.Context, id int) bool {
	res, err := a.db.QueryContext(ctx, drivers.CheckAuthorBYID, id)
	if err != nil || !res.Next() {
		return false
	}

	return true
}
