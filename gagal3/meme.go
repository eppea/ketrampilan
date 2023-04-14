package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

// Transaction struct to represent a transaction
type Transaction struct {
	ID          int     `json:"id"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
}

var db *sql.DB

// Handler function to get all transactions
func getAllTransactions(c echo.Context) error {
	rows, err := db.Query("SELECT id, description, amount FROM transactions")
	if err != nil {
		log.Fatal(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to get transactions",
		})
	}
	defer rows.Close()

	transactions := []Transaction{}
	for rows.Next() {
		var transaction Transaction
		err := rows.Scan(&transaction.ID, &transaction.Description, &transaction.Amount)
		if err != nil {
			log.Fatal(err)
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": "Failed to get transactions",
			})
		}
		transactions = append(transactions, transaction)
	}

	return c.JSON(http.StatusOK, transactions)
}

// Handler function to get a transaction by ID
func getTransaction(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Fatal(err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid transaction ID",
		})
	}

	var transaction Transaction
	err = db.QueryRow("SELECT id, description, amount FROM transactions WHERE id = ?", id).Scan(&transaction.ID, &transaction.Description, &transaction.Amount)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{
				"message": "Transaction not found",
			})
		}
		log.Fatal(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to get transaction",
		})
	}

	return c.JSON(http.StatusOK, transaction)
}

// Handler function to create a new transaction
func createTransaction(c echo.Context) error {
	var transaction Transaction
	err := c.Bind(&transaction)
	if err != nil {
		log.Fatal(err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid request payload",
		})
	}

	result, err := db.Exec("INSERT INTO transactions(description, amount) VALUES (?, ?)", transaction.Description, transaction.Amount)
	if err != nil {
		log.Fatal(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to create transaction",
		})
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to create transaction",
		})
	}

	transaction.ID = int(id)

	return c.JSON(http.StatusCreated, transaction)
}

// Handler function to update a transaction by ID
func updateTransaction(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Fatal(err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid transaction ID",
		})
	}

	var transaction Transaction
	err = c.Bind(&transaction)
	if err != nil {
		log.Fatal(err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid request payload",
		})
	}

	result, err := db.Exec("UPDATE transactions SET description = ?, amount = ? WHERE id = ?", transaction.Description, transaction.Amount, id)
	if err != nil {
		log.Fatal(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to update transaction",
		})
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to update transaction",
		})
	}

	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "Transaction not found",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Transaction updated successfully",
	})
}

// Handler function to delete a transaction by ID
func deleteTransaction(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Fatal(err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid transaction ID",
		})
	}

	result, err := db.Exec("DELETE FROM transactions WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to delete transaction",
		})
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to delete transaction",
		})
	}

	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "Transaction not found",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Transaction deleted successfully",
	})
}

func main() {
	var err error

	// Connect to MySQL database
	db, err = sql.Open("mysql", "root:qwerty@tcp(localhost:3306)/daily_expenses")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// Check if the database is connected
	err = db.Ping()
	if err != nil {
		log.Fatal(err)

	}

	e := echo.New()

	// Define routes
	e.GET("/transactions", getAllTransactions)
	e.GET("/transactions/:id", getTransaction)
	e.POST("/transactions", createTransaction)
	e.PUT("/transactions/:id", updateTransaction)
	e.DELETE("/transactions/:id", deleteTransaction)

	// Start the server
	e.Start(":8080")

}
