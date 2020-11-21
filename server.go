package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Customer struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email`
	Status string `json:"status"`
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

func getAllCustomersHandler(c *gin.Context) {
	customers, err := queryAllCustomer()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}
	c.JSON(http.StatusOK, customers)
}

func getCustomerById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	customer, err := queryByID(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}

	c.JSON(http.StatusOK, customer)
}

func insertNewCustomer(c Customer) (Customer, error) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	}
	defer db.Close()

	row := db.QueryRow("INSERT INTO customers (name,email,status) values($1,$2,$3) RETURNING id", c.Name, c.Email, c.Status)
	var id int
	err = row.Scan(&id)
	if err != nil {
		return Customer{}, fmt.Errorf("can't scan id %s", err)
	}

	customer := Customer{id, c.Name, c.Email, c.Status}
	fmt.Println("insert todo success id: ", id)
	return customer, nil
}

func queryAllCustomer() ([]Customer, error) {
	var customers []Customer
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT id,name,email,status FROM customers")
	if err != nil {
		return []Customer{}, fmt.Errorf("Can't prepare query all customers")
	}

	rows, err := stmt.Query()
	if err != nil {
		return []Customer{}, fmt.Errorf("Can't execute query all customers")
	}

	for rows.Next() {
		var id int
		var name, email, status string

		err := rows.Scan(&id, &name, &email, &status)
		if err != nil {
			return []Customer{}, fmt.Errorf("can't Scan row into variable %s", err)
		}

		customer := Customer{id, name, email, status}
		customers = append(customers, customer)
	}
	fmt.Println("query all customers success")
	return customers, nil
}

func queryByID(rowId int) (Customer, error) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT id, name, email, status FROM customers where id=$1")
	if err != nil {
		return Customer{}, fmt.Errorf("can't prepare query one row statement", err)
	}

	row := stmt.QueryRow(rowId)
	var id int
	var name, email, status string

	err = row.Scan(&id, &name, &email, &status)
	if err != nil {
		return Customer{}, fmt.Errorf("can't scan row into variables %s", err)
	}

	fmt.Println("Customer", id, name, email, status)
	customer := Customer{id, name, email, status}
	return customer, nil

}

func main() {
	r := gin.Default()
	r.GET("/customers", getAllCustomersHandler)
	r.GET("/customers/:id", getCustomerById)
	r.POST("/customers", createCustomerHandler)
	r.Run(":2009")
}
