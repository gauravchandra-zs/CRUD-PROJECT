package datastoreauthor

import (
	"context"
	"database/sql"

	"Projects/GoLang-Interns-2022/threeLayer/drivers"
	"Projects/GoLang-Interns-2022/threeLayer/models"
)

// AuthorStore is a struct
type AuthorStore struct {
	db *sql.DB
}

func New(db *sql.DB) AuthorStore {
	return AuthorStore{
		db,
	}
}

func (a AuthorStore) PostAuthor(ctx context.Context, author models.Author) (int, error) {
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

func (a AuthorStore) PutAuthor(ctx context.Context, id int, author models.Author) (models.Author, error) {
	_, err := a.db.ExecContext(ctx, drivers.UpdateAuthor, author.FirstName, author.LastName, author.Dob, author.PenName, id)
	if err != nil {
		return models.Author{}, err
	}

	return author, nil
}

func (a AuthorStore) DeleteAuthor(ctx context.Context, id int) (int, error) {
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

func (a AuthorStore) GetAuthorByID(ctx context.Context, id int) (models.Author, error) {
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

func (a AuthorStore) CheckAuthor(ctx context.Context, author models.Author) bool {
	row, err := a.db.QueryContext(ctx, drivers.CheckAuthor, author.FirstName, author.LastName, author.Dob, author.PenName)
	if err != nil || !row.Next() {
		return false
	}

	return true
}

func (a AuthorStore) CheckAuthorByID(ctx context.Context, id int) bool {
	res, err := a.db.QueryContext(ctx, drivers.CheckAuthorBYID, id)
	if err != nil || !res.Next() {
		return false
	}

	return true
}
