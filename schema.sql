DROP DATABASE IF EXISTS `transactions_playground`;
CREATE DATABASE `transactions_playground` /*!40100 DEFAULT CHARACTER SET utf8 */;

USE `transactions_playground`;

CREATE TABLE `accounts` (
  `id` int(11) NOT NULL,
  `balance` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
