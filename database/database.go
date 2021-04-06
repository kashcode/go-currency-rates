package database

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

	"example.com/currency-rates/entities"
	_ "github.com/go-sql-driver/mysql"
)

func SaveCurrencyRates(db *sql.DB, rates []entities.CurrencyRate) error {
	log.Println("Saving currency rates into database")

	query := "INSERT INTO rates(currency, rate, date) VALUES "

	var inserts []string
	var params []interface{}

	for _, v := range rates {
		inserts = append(inserts, "(?, ?, ?)")
		params = append(params, v.Currency, v.Rate, v.Date)
	}

	queryVals := strings.Join(inserts, ",")
	query = query + queryVals + " ON DUPLICATE KEY UPDATE rate=VALUES(rate);"

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, params...)
	if err != nil {
		log.Printf("Error %s when inserting row into products table", err)
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d products created simultaneously", rows)

	return nil
}

func getRates(db *sql.DB, query string, args ...interface{}) ([]entities.CurrencyRate, error) {
	rows, err := db.Query(query, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	currencies := []entities.CurrencyRate{}

	for rows.Next() {
		var c entities.CurrencyRate
		if err := rows.Scan(&c.Currency, &c.Rate, &c.Date); err != nil {
			return nil, err
		}
		currencies = append(currencies, c)
	}

	return currencies, nil
}

func GetCurrencyRates(db *sql.DB) ([]entities.CurrencyRate, error) {
	return getRates(db, "SELECT currency, rate, date(`date`) as \"date\" FROM rates where date(`date`) >= date(sysdate());")
}

func GetCurrency(db *sql.DB, currency string) ([]entities.CurrencyRate, error) {
	return getRates(db, "SELECT currency, rate, date(`date`) as \"date\" FROM rates where upper(currency) = upper(?)", currency)
}
