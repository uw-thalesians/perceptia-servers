/*
	Title: Perceptia Database Schema
	Version: 0.1.1
*/
USE master;
GO

IF NOT EXISTS (SELECT * FROM sys.databases WHERE name = 'Perceptia')
	BEGIN
		CREATE DATABASE Perceptia
		COLLATE SQL_Latin1_General_CP1_CI_AS
	END
;
GO

USE Perceptia
;
GO

