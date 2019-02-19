/*
	Title: Perceptia Database Schema
	Version: 0.1.1
*/

-------------------------------------------------------------------------------
-- Setup Database --
-------------------------------------------------------------------------------
-- Select master to remove 
USE [master]
;
GO

-- Remove existing Perceptia DB if applying schema again
If Exists(Select name from master.dbo.sysdatabases Where Name = 'Perceptia')
Begin
	USE [master]
	ALTER DATABASE [Perceptia] SET SINGLE_USER WITH ROLLBACK IMMEDIATE
	DROP DATABASE [Perceptia]
End
;
GO

-- Create Perceptia Database
-- Use same Collate as Azure SQL
CREATE DATABASE Perceptia
	COLLATE Latin1_General_100_CI_AS_SC
;
GO

-------------------------------------------------------------------------------
-- Create Tables --
-------------------------------------------------------------------------------
-- Ensure Perceptia database is selected
USE Perceptia
;
GO

