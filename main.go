package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"database/sql"

	"example.com/currency-rates/bank"
	"example.com/currency-rates/database"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

// Create an exported global variable to hold the database connection pool.
var DB *sql.DB

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbHost := os.Getenv("DATABASE_HOST")
	dbUsername := os.Getenv("DATABASE_USERNAME")
	dbPassword := os.Getenv("DATABASE_PASSWORD")
	dbName := os.Getenv("DATABASE_NAME")
	currencyUrl := os.Getenv("BANK_CCY_URL")

	DB, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", dbUsername, dbPassword, dbHost, dbName))

	if err != nil {
		log.Printf("Error %s when getting db connection", err)
		return
	}
	DB.SetConnMaxLifetime(time.Minute * 3)
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(10)
	defer DB.Close()

	app := &cli.App{
		Name:  "Currency rates",
		Usage: "Reads currency rates from ECB RSS feed and save into DB with API endpoints for displaying rates to user.",
		Commands: []*cli.Command{
			{
				Name:  "bank",
				Usage: "options for task bank",
				Subcommands: []*cli.Command{
					{
						Name:  "get-currency-rates",
						Usage: "get currency rates from RSS feed and store into DB",
						Action: func(c *cli.Context) error {
							bank := &bank.Bank{
								Url: currencyUrl,
							}
							rates, err := bank.GetCurrencyRates()

							err = database.SaveCurrencyRates(DB, rates)
							if err != nil {
								log.Printf("Multiple insert failed with error %s", err)
								return err
							}

							return nil
						},
					},
				},
			},
			{
				Name:  "api",
				Usage: " run currency rates endpoint API",
				Action: func(c *cli.Context) error {
					r := mux.NewRouter()
					r.HandleFunc("/currencies", func(w http.ResponseWriter, r *http.Request) {
						w.Header().Set("Content-Type", "application/json")

						currencies, err := database.GetCurrencyRates(DB)
						if err != nil {
							w.WriteHeader(http.StatusInternalServerError)
							json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
							return
						}

						response, _ := json.Marshal(currencies)
						w.WriteHeader(http.StatusOK)
						w.Write(response)
					})
					r.HandleFunc("/currencies/{ccy}", func(w http.ResponseWriter, r *http.Request) {
						vars := mux.Vars(r)
						w.Header().Set("Content-Type", "application/json")

						currencies, err := database.GetCurrency(DB, vars["ccy"])
						if err != nil {
							w.WriteHeader(http.StatusInternalServerError)
							json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
							return
						}

						response, _ := json.Marshal(currencies)
						w.WriteHeader(http.StatusOK)
						w.Write(response)
					})

					log.Fatal(http.ListenAndServe(":8000", r))

					return nil
				},
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
