CREATE DATABASE if not exists whim_development;
use whim_development;


CREATE TABLE if not exists `conversions` (
  `id` bigint(20) unsigned NOT NULL PRIMARY KEY AUTO_INCREMENT,
  `currency_id_from` bigint(20) NOT NULL,
  `currency_id_to` bigint(20) NOT NULL,
  `rate` float NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE if not exists `currencies` (
  `id` bigint(20) unsigned NOT NULL PRIMARY KEY AUTO_INCREMENT,
  `name` varchar(50) NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
