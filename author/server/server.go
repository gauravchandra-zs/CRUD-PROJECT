package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Author struct {
	ID        int    `json:",omitempty"`
	FirstName string `json:",omitempty"`
	LastName  string `json:",omitempty"`
	Dob       string `json:",omitempty"`
	PenName   string `json:",omitempty"`
}
type Book struct {
	ID              int    `json:",omitempty"`
	Title           string `json:",omitempty"`
	Author          Author `json:",omitempty"`
	Publication     string `json:",omitempty"`
	PublicationDate string `json:",omitempty"`
}

// const  are constants used in methods
const (
	Driver = "mysql"

	DataSource = "root:@Gc18072001@/test"

	CreateBook = "CREATE Table IF NOT EXISTS Book(bookID int NOT NULL AUTO_INCREMENT, title varchar(50), " +
		"publication varchar(50),publicationDate varchar(50),authorID int, PRIMARY KEY (bookID),FOREIGN KEY " +
		"(authorID) REFERENCES Author(authorID));"

	InsertIntoBook = "INSERT INTO Book (title,publication,publicationDate,authorID) " +
		"VALUES (? ,?, ?, ?)"
	DeleteBookQuery = "DELETE FROM Book WHERE BookID = ?"
	CheckBook       = "SELECT bookID FROM Book WHERE title = ? AND publication =? AND" +
		"  publicationDate= ?AND authorID =?"
	CheckBookBYID         = "SELECT * FROM Book WHERE bookID = ?"
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

// createTable create table if table not exists
func createTable() (*sql.DB, error) {
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

// checkBook check book if exist then return true otherwise false
func checkBook(book *Book, authorID int64) bool {
	db, err := sql.Open(Driver, DataSource)
	if err != nil {
		log.Print(err)
		return false
	}

	result, err := db.Query(CheckBook, book.Title, book.Publication, book.PublicationDate, authorID)

	if err != nil || !result.Next() {
		return false
	}

	return true
}

// checkAuthor check author if exist then return true otherwise false
func checkAuthor(author Author) bool {
	db, err := sql.Open(Driver, DataSource)
	if err != nil {
		log.Print(err)
		return false
	}

	res, err := db.Query(CheckAuthor, author.FirstName, author.LastName, author.Dob, author.PenName)
	if !res.Next() || err != nil {
		return false
	}

	return true
}
func checkAuthorByID(id int) bool {
	db, err := sql.Open(Driver, DataSource)
	if err != nil {
		log.Print(err)
		return false
	}

	res, err := db.Query(CheckAuthorBYID, id)
	if !res.Next() || err != nil {
		return false
	}

	return true
}

func checkBookBid(id int) bool {
	db, err := sql.Open(Driver, DataSource)
	if err != nil {
		log.Print(err)
		return false
	}

	res, err := db.Query(CheckBookBYID, id)
	if !res.Next() || err != nil {
		return false
	}

	return true
}

func ValidateBook(b *Book) bool {
	slc := strings.Split(b.PublicationDate, "-")
	sz := 3

	switch {
	case b.Publication != "Scholastic" && b.Publication != "Penguin" && b.Publication != "Arihanth":
		return false
	case len(slc) < sz:
		return false
	case slc[2] >= "2022" || slc[2] < "1880":
		return false
	case b.Title == "":
		return false
	default:
		return true
	}
}

func ValidateAuthor(author Author) bool {
	if author.FirstName == "" || author.LastName == "" || author.Dob == "" || author.PenName == "" {
		return false
	}

	return true
}

func GetAuthorByID(id int) (Author, error) {
	db, err := createTable()
	if err != nil {
		return Author{}, err
	}

	ResAuthor, err2 := db.Query(SelectAuthorByID, id)
	if err2 != nil {
		return Author{}, err
	}

	author := Author{}

	if ResAuthor.Next() {
		err2 = ResAuthor.Scan(&author.ID, &author.FirstName, &author.LastName, &author.Dob, &author.PenName)
		if err2 != nil {
			return Author{}, err
		}
	}

	return author, nil
}

func GetID(req *http.Request) (int, error) {
	params := mux.Vars(req)
	ID, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Print("error in converting string id to int")
	}

	if ID <= 0 {
		return 0, errors.New("id not found")
	}

	return ID, nil
}

func GetBookByID(id int) (Book, int, error) {
	db, err := createTable()
	if err != nil {
		return Book{}, 0, err
	}

	result, err := db.Query(SelectFromBookByID, id)
	if err != nil {
		return Book{}, 0, err
	}

	if result.Next() {
		var authorID int

		book := Book{}

		err = result.Scan(&book.ID, &book.Title, &book.Publication, &book.PublicationDate, &authorID)
		if err != nil {
			return Book{}, 0, err
		}

		return book, authorID, nil
	}

	return Book{}, 0, errors.New("book not exists")
}
func ReadAuthor(req *http.Request, author *Author) error {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &author)
	if err != nil {
		return err
	}

	return nil
}
func ReadBook(req *http.Request, book *Book) error {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &book)
	if err != nil {
		return err
	}

	return nil
}

func PostBookHelper(author Author) int64 {
	db, err := createTable()
	if err != nil {
		return 0
	}

	var authorID int64

	result, err := db.Query(CheckAuthor, author.FirstName, author.LastName, author.Dob, author.PenName)
	if err != nil {
		return 0
	}

	if result.Next() {
		err = result.Scan(&authorID)
		if err != nil {
			return 0
		}
	} else {
		res, err3 := db.Exec(InsertIntoAuthor, author.FirstName, author.LastName, author.Dob, author.PenName)
		if err3 != nil {
			return 0
		}
		authorID, err3 = res.LastInsertId()
		if err3 != nil {
			return 0
		}
	}

	return authorID
}
func GetAllHelper(req *http.Request) (string, *sql.Rows, error) {
	db, err := createTable()
	if err != nil {
		return "", nil, err
	}

	title := req.URL.Query().Get("title")
	authorInclude := req.URL.Query().Get("includeAuthor")

	var result *sql.Rows
	if title == "" {
		result, err = db.Query(SelectFromBook)
		if err != nil {
			return "", result, err
		}
	} else {
		result, err = db.Query(SelectFromBookByTitle, title)
		if err != nil {
			return "", result, err
		}
	}

	return authorInclude, result, nil
}
func GetAll(w http.ResponseWriter, req *http.Request) {
	authorInclude, res, err := GetAllHelper(req)
	if err != nil {
		return
	}

	var books []Book

	var id int

	b := Book{}

	for res.Next() {
		err = res.Scan(&b.ID, &b.Title, &b.Publication, &b.PublicationDate, &id)
		if err != nil {
			return
		}

		if authorInclude == "true" {
			b.Author, err = GetAuthorByID(id)
			if err != nil {
				return
			}
		}

		books = append(books, b)
	}

	body, _ := json.Marshal(books)
	_, err = w.Write(body)

	if err != nil {
		return
	}
}
func GetByID(w http.ResponseWriter, req *http.Request) {
	ID, err := GetID(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	book, authorID, err := GetBookByID(ID)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	author, err := GetAuthorByID(authorID)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	book.Author = author

	data, err := json.Marshal(book)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(data)
	if err != nil {
		return
	}
}
func PostBook(w http.ResponseWriter, req *http.Request) {
	db, err := createTable()
	if err != nil {
		return
	}

	var book Book

	err = ReadBook(req, &book)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !ValidateBook(&book) || !ValidateAuthor(book.Author) {
		fmt.Println(ValidateBook(&book), ValidateAuthor(book.Author))
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	authorID := PostBookHelper(book.Author)
	if !checkBook(&book, authorID) {
		_, err = db.Exec(InsertIntoBook, book.Title, book.Publication, book.PublicationDate, authorID)
		if err != nil {
			return
		}

		w.WriteHeader(http.StatusCreated)
	}

	w.WriteHeader(http.StatusCreated)
}
func PostAuthor(w http.ResponseWriter, req *http.Request) {
	db, err := createTable()
	if err != nil {
		log.Print(err)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Print(err)
		return
	}

	var author Author
	err = json.Unmarshal(body, &author)

	if err != nil {
		log.Print(err)
		return
	}

	if !ValidateAuthor(author) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !checkAuthor(author) {
		_, err = db.Exec(InsertIntoAuthor, author.FirstName, author.LastName, author.Dob, author.PenName)
		if err != nil {
			return
		}

		w.WriteHeader(http.StatusCreated)

		return
	}

	w.WriteHeader(http.StatusCreated)
}
func DeleteBook(w http.ResponseWriter, req *http.Request) {
	db, err := createTable()
	if err != nil {
		log.Print(err)
		return
	}

	params := mux.Vars(req)

	ID, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Print(err)
	}

	if ID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !checkBookBid(ID) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = db.Exec(DeleteBookQuery, ID)
	if err != nil {
		log.Print("book not exist")
	}
}
func DeleteAuthor(w http.ResponseWriter, req *http.Request) {
	db, err := createTable()
	if err != nil {
		return
	}

	ID, err := GetID(req)
	if ID <= 0 || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !checkAuthorByID(ID) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM Book WHERE AuthorID=?", ID)
	if err != nil {
		return
	}

	_, err = db.Exec(DeleteAuthorQuery, ID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	d, _ := json.Marshal(ID)

	_, err = w.Write(d)
	if err != nil {
		return
	}
}

// PutAuthor update author detail with given id
func PutAuthor(w http.ResponseWriter, req *http.Request) {
	db, err := createTable()
	if err != nil {
		return
	}

	var author Author
	err = ReadAuthor(req, &author)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !ValidateAuthor(author) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ID, err := GetID(req)
	if ID <= 0 || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !checkAuthorByID(ID) {
		log.Print("id not present")
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	_, err = db.Exec(UpdateAuthor, author.FirstName, author.LastName, author.Dob, author.PenName, ID)
	if err != nil {
		log.Print(err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// PutBook  update book detail with given id
func PutBook(w http.ResponseWriter, req *http.Request) {
	db, err := createTable()
	if err != nil {
		return
	}

	var book Book

	err = ReadBook(req, &book)
	if err != nil || !ValidateBook(&book) || !ValidateAuthor(book.Author) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ID, err := GetID(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !checkAuthorByID(book.Author.ID) || !checkBookBid(book.ID) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = db.Exec(UpdateBook, book.Title, book.Publication, book.PublicationDate, ID)
	if err != nil {
		return
	}

	_, err = db.Exec(UpdateAuthor, book.Author.FirstName, book.Author.LastName, book.Author.Dob, book.Author.PenName, book.Author.ID)
	if err != nil {
		log.Print(err)
	}

	w.WriteHeader(http.StatusAccepted)
}
func main() {
	r := mux.NewRouter()

	r.HandleFunc("/books", GetAll).Methods(http.MethodGet)                // done
	r.HandleFunc("/books/{id}", GetByID).Methods(http.MethodGet)          // done
	r.HandleFunc("/books", PostBook).Methods(http.MethodPost)             // done
	r.HandleFunc("/author", PostAuthor).Methods(http.MethodPost)          // done
	r.HandleFunc("/author/{id}", DeleteAuthor).Methods(http.MethodDelete) // done
	r.HandleFunc("/books/{id}", DeleteBook).Methods(http.MethodDelete)    // done
	r.HandleFunc("/author/{id}", PutAuthor).Methods(http.MethodPut)       // done
	r.HandleFunc("/books/{id}", PutBook).Methods(http.MethodPut)          // done

	server := http.Server{
		Addr:    ":8000",
		Handler: r,
	}

	fmt.Println("server stared")

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
