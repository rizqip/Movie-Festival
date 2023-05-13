-- Adminer 4.8.1 MySQL 10.11.2-MariaDB dump

SET NAMES utf8;
SET time_zone = '+00:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

SET NAMES utf8mb4;

CREATE DATABASE `MovieFestival` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci */;
USE `MovieFestival`;

CREATE TABLE `Artists` (
  `ArtistId` bigint(20) NOT NULL AUTO_INCREMENT,
  `ArtistName` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`ArtistId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


CREATE TABLE `Genres` (
  `GenreId` bigint(20) NOT NULL AUTO_INCREMENT,
  `GenreName` varchar(50) DEFAULT NULL,
  `ViewCount` int(11) DEFAULT 0,
  PRIMARY KEY (`GenreId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


CREATE TABLE `Movies` (
  `MovieId` bigint(20) NOT NULL AUTO_INCREMENT,
  `Title` varchar(255) DEFAULT NULL,
  `Description` longtext DEFAULT NULL,
  `Duration` varchar(50) DEFAULT NULL,
  `ArtistId` varchar(50) DEFAULT NULL,
  `GenreId` varchar(50) DEFAULT NULL,
  `Url` text DEFAULT NULL,
  `ViewCount` int(11) DEFAULT 0,
  `CreatedAt` datetime DEFAULT NULL,
  `CreatedBy` bigint(20) DEFAULT NULL,
  `UpdatedAt` datetime DEFAULT NULL,
  PRIMARY KEY (`MovieId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


CREATE TABLE `MovieVotes` (
  `VoteId` bigint(20) NOT NULL AUTO_INCREMENT,
  `MovieId` bigint(20) DEFAULT NULL,
  `UserId` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`VoteId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


CREATE TABLE `UserRecords` (
  `RecordId` bigint(20) NOT NULL AUTO_INCREMENT,
  `UserId` bigint(20) DEFAULT NULL,
  `MovieId` bigint(20) DEFAULT NULL,
  `CreatedAt` datetime DEFAULT NULL,
  PRIMARY KEY (`RecordId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


CREATE TABLE `Users` (
  `UserId` bigint(20) NOT NULL AUTO_INCREMENT,
  `Name` varchar(255) DEFAULT NULL,
  `Email` varchar(100) DEFAULT NULL,
  `Password` varchar(50) DEFAULT NULL,
  `Type` tinyint(4) DEFAULT NULL COMMENT '1. User, 2. Admin',
  `Status` tinyint(4) DEFAULT NULL COMMENT '0. Logout, 1. Login',
  `CreatedAt` datetime DEFAULT NULL,
  `UpdatedAt` datetime DEFAULT NULL,
  PRIMARY KEY (`UserId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


-- 2023-05-13 05:09:07
