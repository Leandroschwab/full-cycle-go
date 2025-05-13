package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"github.com/devfullcycle/20-CleanArch/internal/db"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	ctx := context.Background()
	//dbConn, err := sql.Open("mysql", "root:root@tcp(172.20.20.15:3306)/orders")
	server := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	user := os.Getenv("MYSQL_USER")	
	password := os.Getenv("MYSQL_PASSWORD")
	database := os.Getenv("MYSQL_DATABASE")
	// Connect to the database using the provided credentials
	dbConn, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, server, port, database))
	if err != nil {
		panic(err)
	}
	defer dbConn.Close()

	queries := db.New(dbConn)

	// Check if the table exists before creating it
	if !tableExists(ctx, dbConn, "orders") {
		err = queries.CreateOrderTable(ctx)
		if err != nil {
			panic(err)
		}
		fmt.Println("Table 'orders' created successfully.")
	} else {
		fmt.Println("Table 'orders' already exists.")
	}

	// Check if the order already exists before inserting
	orderID := "1"
	if !orderExists(ctx, queries, orderID) {
		err = queries.CreateOrder(ctx, db.CreateOrderParams{
			ID:         orderID,
			Price:      100,
			Tax:        10,
			FinalPrice: 110,
		})
		if err != nil {
			panic(err)
		}
		fmt.Printf("Order with ID '%s' created successfully.\n", orderID)
	} else {
		fmt.Printf("Order with ID '%s' already exists.\n", orderID)
	}

	orders, err := queries.GetAllOrders(ctx)
	if err != nil {
		panic(err)
	}
	for _, order := range orders {
		fmt.Printf("Order: ID=%s, Price=%.2f, Tax=%.2f, FinalPrice=%.2f\n", order.ID, order.Price, order.Tax, order.FinalPrice)
	}
}

func tableExists(ctx context.Context, dbConn *sql.DB, tableName string) bool {
	query := fmt.Sprintf("SHOW TABLES LIKE '%s'", tableName)
	var result string
	err := dbConn.QueryRowContext(ctx, query).Scan(&result)
	return err == nil
}

func orderExists(ctx context.Context, queries *db.Queries, orderID string) bool {
	_, err := queries.GetOrderById(ctx, orderID)
	return err == nil
}
