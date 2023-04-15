package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

// struntnya
type Transaction struct {
	ID          int     `json:"id"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
}

var db *sql.DB

// semua transaksi
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

// transaksi per ID
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

// transaksi baru
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

// transaksi update
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

// delet
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

	// konek database
	db, err = sql.Open("mysql", "root:EANHHUFWsX2ocayI4WXW@tcp(containers-us-west-143.railway.app:6475)/railway")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// jika konek
	err = db.Ping()
	if err != nil {
		log.Fatal(err)

	}

	e := echo.New()

	// routes
	e.GET("/transactions", getAllTransactions)
	e.GET("/transactions/:id", getTransaction)
	e.POST("/transactions", createTransaction)
	e.PUT("/transactions/:id", updateTransaction)
	e.DELETE("/transactions/:id", deleteTransaction)

	// start
	e.Start(":" + getPort())

}
func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Port default jika tidak ada environment variable
	}
	return port
}
