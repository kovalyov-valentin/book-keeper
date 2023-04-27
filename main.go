package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Person struct {
	gorm.Model

	Name  string
	Email string `gorm:"typevarchar(100);unique_index"`
	Books []Book
}

type Book struct {
	gorm.Model

	Title      string
	Author     string
	CallNumber int `gorm:"unique_index"`
	PersonID   int
}

// var (
// 	person = &Person{
// 		Name:  "Valentin",
// 		Email: "valentin2011k1997@maol.ru",
// 	}
// 	books = []Book{
// 		{Title: "Microservies",
// 			Author:     "Kris Richards",
// 			CallNumber: 123,
// 			PersonID:   1,
// 		},
// 		{Title: "Head First",
// 			Author:     "Эрик Фримен",
// 			CallNumber: 234,
// 			PersonID:   1,
// 		},
// 	}
// )

var db *gorm.DB
var err error

func main() {
	// Loading environment variables
	dialect := os.Getenv("DIALECT")
	host := os.Getenv("HOST")
	dbPort := os.Getenv("DBPORT")
	user := os.Getenv("USER")
	dbName := os.Getenv("NAME")
	dbpassword := os.Getenv("PASSWORD")

	//Database connection string

	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, user, dbName, dbpassword, dbPort)

	// Openning connection
	db, err = gorm.Open(dialect, dbURI)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Successfully connected to database!")
	}

	// Close connection to database when the main function finishes
	defer db.Close()

	//Make migrations to the database if they hane not already been created
	db.AutoMigrate(&Person{})
	db.AutoMigrate(&Book{})

	// db.Create(&person)
	// for idx := range books {
	// 	db.Create(&books[idx])
	// }

	// API routes 
	router := mux.NewRouter()

	router.HandleFunc("/people", getPeople).Methods("GET")
	router.HandleFunc("/person/{id}", getPerson).Methods("GET") // and their books
	router.HandleFunc("/create/person", createPerson).Methods("POST")
	router.HandleFunc("/delete/person/{id}", deletePerson).Methods("DELETE")
	// router.HandleFunc("/update/person/[id]", updatePerson).Methods("PUT")

	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/book/{id}", getBook).Methods("GET")
	router.HandleFunc("/create/book", createBook).Methods("POST")
	router.HandleFunc("/delete/book/{id}", deleteBook).Methods("DELETE")
	// router.HandleFunc("/update/book/{id}", updateBook).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8080", router))

}
// API Controllers
func getPeople(w http.ResponseWriter, r *http.Request) {
	var people []Person 

	db.Find(&people)

	json.NewEncoder(w).Encode(&people)


}

func getPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var person Person 
	var books []Book

	db.First(&person, params["id"])
	db.Model(&person).Related(&books)
	
	person.Books = books

	json.NewEncoder(w).Encode(person)
}

func createPerson(w http.ResponseWriter, r *http.Request) {
	var person Person

	json.NewDecoder(r.Body).Decode(&person)

	createPerson := db.Create(&person)
	err = createPerson.Error 
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(&person)
	}
}

func deletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var person Person 

	db.First(&person, params["id"])
	db.Delete(&person)

	json.NewEncoder(w).Encode(&person)

}

func getBooks(w http.ResponseWriter, r *http.Request) {
	var books []Book 

	db.Find(&books)

	json.NewEncoder(w).Encode(&books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var book Book 

	db.First(&book, params["id"])

	json.NewEncoder(w).Encode(book)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	var book Book 

	json.NewDecoder(r.Body).Decode(&book)

	createBook := db.Create(&book)
	err = createBook.Error
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(&book)
	}

}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var book Book 

	db.First(&book, params["id"])
	db.Delete(&book)

	json.NewEncoder(w).Encode(&book)
}

