/*
	Title: Perceptia Database Schema
	Version: 0.2.0
*/
-------------------------------------------------------------------------------
-- Change Log --
-------------------------------------------------------------------------------
/*
	Date, Changer, Short Description, Version
	2019/02/19, Chris, Created Schema, 0.1.0
	2019/02/19, Chris, Change DB Collation to support _SC, 0.1.1
	2019/02/19, Chris, Add Table Definitions, 0.2.0
*/

-------------------------------------------------------------------------------
-- Sections --
-------------------------------------------------------------------------------

-- Setup Database
-- Create Tables
-- Create Indexes

-------------------------------------------------------------------------------
-- Setup Database --
-------------------------------------------------------------------------------
-- Select master to remove 
USE [master]
;
GO

-- Remove existing Perceptia DB if applying schema again
If Exists(SELECT [name] FROM master.dbo.sysdatabases WHERE [name] = 'Perceptia')
Begin
	USE [master]
	ALTER DATABASE [Perceptia] SET SINGLE_USER WITH ROLLBACK IMMEDIATE
	DROP DATABASE [Perceptia]
End
;
GO

-- Create Perceptia Database
-- Use same Collate as Azure SQL
CREATE DATABASE [Perceptia]
	COLLATE Latin1_General_100_CI_AS_SC
;
GO

-------------------------------------------------------------------------------
-- Create Tables --
-------------------------------------------------------------------------------
-- Ensure Perceptia database is selected
USE [Perceptia]
;
GO

-----------------------------------------------------------
-- Account Table --
-----------------------------------------------------------
-- Summary: Store information about a specific user

CREATE TABLE [Account] (
	[AccountUUID] UNIQUEIDENTIFIER DEFAULT(NEWID()) NOT NULL,
	[Username] NVARCHAR(255) NOT NULL,
	[FullName] NVARCHAR(255),
	[DisplayName] NVARCHAR(255),
	[Email] NVARCHAR(255) NOT NULL,
	[Created] DATETIME DEFAULT(GETDATE()),
	[EncodedPassword] NVARCHAR(500) NOT NULL,
	CONSTRAINT [PK_Account_AccountUUID] PRIMARY KEY ([AccountUUID]),
	CONSTRAINT [UQ_Account_Username] UNIQUE ([Username]),
	CONSTRAINT [UQ_Account_Email] UNIQUE ([Email])
)
;
GO


-------------------------------------------------------------------------------
-- Create Indexes --
-------------------------------------------------------------------------------

-----------------------------------------------------------
-- Account Table --
-----------------------------------------------------------

-- Index for the Username column
CREATE INDEX [IX_Account_Username]
	ON [Account] ([Username])
;
GO

-- Index for the Email column
CREATE INDEX [IX_Account_Email]
	ON [Account] ([Email])
;
GO

