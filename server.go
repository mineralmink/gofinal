package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Customer struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email`
	Status string `json:"status"`
}

func helloHandler(w http.ResponseWriter, req *http.Request) {
	resp := []byte(`{"name": "anuchit"}`)
	w.Write(resp)
}
func createDatabase() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	}
	defer db.Close()

	createTb := `
	CREATE TABLE IF NOT EXISTS customers (
		id SERIAL PRIMARY KEY,
		name TEXT,
		email TEXT,
		status TEXT
	);
	`
	_, err = db.Exec(createTb)

	if err != nil {
		log.Fatal("can't create table", err)
	}

	fmt.Println("create table success")

	fmt.Println("okay")
}

func createCustomerHandler(c *gin.Context) {
	customer := Customer{}

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cus, err := insertNewCustomer(customer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, cus)
}

func insertNewCustomer(c Customer) ([]Customer, error) {
	var customers []Customer
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	}
	defer db.Close()

	row := db.QueryRow("INSERT INTO customers (name,email,status) values($1,$2,$3) RETURNING id", c.Name, c.Email, c.Status)
	var id int
	err = row.Scan(&id)
	if err != nil {
		return []Customer{}, fmt.Errorf("can't scan id %s", err)
	}

	customer := Customer{id, c.Name, c.Email, c.Status}
	customers = append(customers, customer)
	fmt.Println("insert todo success id: ", id)
	return customers, nil

}

func main() {
	r := gin.Default()
	r.POST("/customers", createCustomerHandler)
	r.Run(":2009")
}
