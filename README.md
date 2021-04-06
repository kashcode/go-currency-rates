# go-currency-rates

docker run -d -p 3306:3306 --name mariadb -e MYSQL_ROOT_PASSWORD=secret -v /mysql mariadb

docker container start mariadb

CREATE SCHEMA `currency` DEFAULT CHARACTER SET utf8mb4 ;

CREATE TABLE `currency`.`rates` (
  `currency` VARCHAR(3) NOT NULL,
  `rate` DECIMAL(19,4) NOT NULL,
  `date` DATETIME NOT NULL,
  UNIQUE INDEX `ccy_date` (`currency` ASC, `date` ASC));