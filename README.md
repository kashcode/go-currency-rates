# go-currency-rates

Configuration parameters stored in `.env` file.

1. run `docker-compose up -d` - this command will start MariaDB

2. connect to DB and run sql script:
```sql
CREATE SCHEMA `currency` DEFAULT CHARACTER SET utf8mb4;

CREATE TABLE `currency`.`rates` (
  `currency` VARCHAR(3) NOT NULL,
  `rate` DECIMAL(19,4) NOT NULL,
  `date` DATETIME NOT NULL,
  UNIQUE INDEX `ccy_date` (`currency` ASC, `date` ASC));
```

3. run `go run . bank get-currency-rates` - to get and save currency rates
4. run `go run . api` - to start API endpoints
 
   4.1 `http://127.0.0.1:8000/currencies` - return in JSON format latest curency rates.

   ```json
   [{"currency":"AUD","rate":1.5482,"date":"2021-04-06T00:00:00Z"},{"currency":"BGN","rate":1.9558, ...
   ```
   4.2 `http://127.0.0.1:8000/currencies/{ccy}` - historical values for a particular currency rate .
   
   **ccy** - currency code. `AUD`, `USD` ... etc.
   
   Example:

   `http://127.0.0.1:8000/currencies/usd`:

   ```json
   [{"currency":"USD","rate":1.1812,"date":"2021-04-06T00:00:00Z"}]
   ```