package main

import (
	_ "net/http"
	"os"

	"developer.zopsmart.com/go/gofr/pkg/gofr"

	_ "github.com/go-sql-driver/mysql"

	datastoreauthor "MyProject/CRUD-PROJECT/threeLayer/datastore/author"
	datastorebook "MyProject/CRUD-PROJECT/threeLayer/datastore/book"
	handlerauthor "MyProject/CRUD-PROJECT/threeLayer/handler/author"
	handlerbook "MyProject/CRUD-PROJECT/threeLayer/handler/book"
	serviceauthor "MyProject/CRUD-PROJECT/threeLayer/service/author"
	servicebook "MyProject/CRUD-PROJECT/threeLayer/service/book"
)

func main() {
	//db, err := drivers.CreateTable()
	//if err != nil {
	//	log.Print(err)
	//	return
	//}
	err := os.Setenv("DB_HOST", "localhost")
	if err != nil {
		return
	}

	err = os.Setenv("DB_USER", "root")
	if err != nil {
		return
	}

	err = os.Setenv("DB_PASSWORD", "gaurav")
	if err != nil {
		return
	}

	err = os.Setenv("DB_NAME", "test")
	if err != nil {
		return
	}

	err = os.Setenv("DB_PORT", "1000")
	if err != nil {
		return
	}

	err = os.Setenv("DB_DIALECT", "mysql")
	if err != nil {
		return
	}

	err = os.Setenv("APP_NAME", "CRUD")
	if err != nil {
		return
	}

	err = os.Setenv("APP_VERSION", "0.1")
	if err != nil {
		return
	}

	err = os.Setenv("HTTP_PORT", "8000")
	if err != nil {
		return
	}

	bookStore := datastorebook.New()
	authorStore := datastoreauthor.New()

	svcBook := servicebook.New(bookStore, authorStore)
	svcAuthor := serviceauthor.New(authorStore, bookStore)

	book := handlerbook.New(svcBook)
	author := handlerauthor.New(svcAuthor)

	r := gofr.New()

	r.Server.ValidateHeaders = false

	r.GET("/book", book.GetAllBooks)
	r.GET("/book/{id}", book.GetBookByID)
	r.POST("/book", book.PostBook)
	r.PUT("/book/{id}", book.PutBook)
	r.DELETE("/book/{id}", book.DeleteBook)

	r.POST("/author", author.PostAuthor)
	r.DELETE("/author/{id}", author.DeleteAuthor)
	r.PUT("/author/{id}", author.PutAuthor)
	r.Start()
}
