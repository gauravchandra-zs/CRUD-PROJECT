package drivers

import (
	"database/sql"
)

const (
	Driver = "mysql"

	//	DataSource = "root:@Gc18072001@/test"

	DataSource = "root:gaurav@tcp(localhost:1000)/test"

	CreateBook = "CREATE Table IF NOT EXISTS Book(bookID int NOT NULL AUTO_INCREMENT, title varchar(50), " +
		"publication varchar(50),publicationDate varchar(50),authorID int, PRIMARY KEY (bookID),FOREIGN KEY " +
		"(authorID) REFERENCES Author(authorID));"

	InsertIntoBook = "INSERT INTO Book (title,publication,publicationDate,authorID) " +
		"VALUES (? ,?, ?, ?)"
	DeleteBookQuery       = "DELETE FROM Book WHERE BookID = ?"
	DeleteBookByAuthorID  = "DELETE FROM Book WHERE AuthorID=?"
	CheckBook             = "SELECT bookID FROM Book WHERE title = ? AND publication =? AND publicationDate= ? AND authorID =?"
	CheckBookBYID         = "SELECT bookID FROM Book WHERE bookID = ?"
	SelectFromBook        = "SELECT * FROM Book;"
	SelectFromBookByTitle = "SELECT * FROM Book WHERE title =?;"
	SelectFromBookByID    = "SELECT * FROM Book WHERE bookID= ?"
	UpdateBook            = "UPDATE Book SET title = ? ,publication = ? " +
		",publicationDate = ?  WHERE bookID =?"

	CreateAuthor = "CREATE Table IF NOT EXISTS Author(authorID int NOT NULL AUTO_INCREMENT, " +
		"firstName varchar(50), lastName varchar(50),dob varchar(50),penName varchar(50), PRIMARY KEY (authorID));"
	InsertIntoAuthor = "INSERT INTO Author (firstName,lastName,dob,penName) VALUES (? ,?, ?, ?)"
	CheckAuthor      = "SELECT authorID FROM Author WHERE firstName = ? AND " +
		"lastName =? AND  dob = ? AND penName =?"
	CheckAuthorBYID   = "SELECT authorID FROM Author WHERE authorID = ?"
	DeleteAuthorQuery = "DELETE  FROM Author WHERE authorID = ?"
	SelectAuthorByID  = "SELECT * FROM Author WHERE authorID=?"
	UpdateAuthor      = "UPDATE Author SET firstName = ? ,lastName = ? ,dob = ? ,penName = ?  WHERE authorID =?"
)

func CreateTable() (*sql.DB, error) {
	db, err := sql.Open(Driver, DataSource)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(CreateAuthor)

	if err != nil {
		return db, err
	}

	_, err = db.Exec(CreateBook)

	if err != nil {
		return db, err
	}

	return db, nil
}
