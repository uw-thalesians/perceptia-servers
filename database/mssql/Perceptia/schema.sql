/*
	Title: Perceptia Database Schema
	Version: 0.3.0
*/
-------------------------------------------------------------------------------
-- Change Log --
-------------------------------------------------------------------------------
/*
	Date, Changer, Short Description, Version
	2019/02/19, Chris, Created Schema, 0.1.0
	2019/02/19, Chris, Change DB Collation to support _SC, 0.1.1
	2019/02/19, Chris, Add Table Definitions, 0.2.0
	2019/02/28, Chris, Change table structure, 0.3.0
*/

-------------------------------------------------------------------------------
-- Sections --
-------------------------------------------------------------------------------

-- Setup Database
-- Create Tables
-- Create Foreign Key Constraints
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
-- User Table --
-----------------------------------------------------------
-- Summary: Store information about a specific user

CREATE TABLE [User] (
	[UUID] UNIQUEIDENTIFIER DEFAULT(NEWID()) NOT NULL,
	[Username] NVARCHAR(255) NOT NULL,
	[FullName] NVARCHAR(255),
	[DisplayName] NVARCHAR(255),
	[Created] DATETIME DEFAULT(GETDATE()),
	CONSTRAINT [PK_User_UUID] PRIMARY KEY ([UUID]),
	CONSTRAINT [UQ_User_Username] UNIQUE ([Username])
)
;
GO

-----------------------------------------------------------
-- Email Table --
-----------------------------------------------------------
-- Summary: Store email addresses

CREATE TABLE [Email] (
	[UUID] UNIQUEIDENTIFIER DEFAULT(NEWID()) NOT NULL,
	[Email] NVARCHAR(255) NOT NULL,
	[Created] DATETIME DEFAULT(GETDATE()),
	CONSTRAINT [PK_Email_UUID] PRIMARY KEY ([UUID])
)
;
GO

-----------------------------------------------------------
-- UserEmail Table --
-----------------------------------------------------------
-- Summary: Associates an email with a user

CREATE TABLE [UserEmail] (
	[UUID] UNIQUEIDENTIFIER DEFAULT(NEWID()) NOT NULL,
	[User_UUID] UNIQUEIDENTIFIER NOT NULL,
	[Email_UUID] UNIQUEIDENTIFIER NOT NULL,
	CONSTRAINT [PK_UserEmail_UUID] PRIMARY KEY ([UUID]),
)
;
GO

-----------------------------------------------------------
-- Credential Table --
-----------------------------------------------------------
-- Summary: Store login credentials

CREATE TABLE [Credential] (
	[UUID] UNIQUEIDENTIFIER DEFAULT(NEWID()) NOT NULL,
	[Created] DATETIME DEFAULT(GETDATE()),
	[EncodedHash] NVARCHAR(500) NOT NULL,
	CONSTRAINT [PK_Credential_UUID] PRIMARY KEY ([UUID])
)
;
GO

-----------------------------------------------------------
-- UserCredential Table --
-----------------------------------------------------------
-- Summary: Associates a user with a login credential

CREATE TABLE [UserCredential] (
	[UUID] UNIQUEIDENTIFIER DEFAULT(NEWID()) NOT NULL,
	[User_UUID] UNIQUEIDENTIFIER NOT NULL,
	[Credential_UUID] UNIQUEIDENTIFIER NOT NULL,
	CONSTRAINT [PK_UserCredential_UUID] PRIMARY KEY ([UUID]),
	CONSTRAINT [UQ_UserCredential_UserUUID] UNIQUE ([User_UUID]),
	CONSTRAINT [UQ_UserCredential_CredentialUUID] UNIQUE ([Credential_UUID])
)
;
GO


-------------------------------------------------------------------------------
-- Create Foreign Key Constraints --
-------------------------------------------------------------------------------

-----------------------------------------------------------
-- UserEmail Table --
-----------------------------------------------------------

ALTER TABLE [UserEmail]
	ADD
	CONSTRAINT [FK_UserEmail_UserUUID] FOREIGN KEY ([User_UUID])
		REFERENCES [User] ([UUID]),
	CONSTRAINT [FK_UserEmail_EmailUUID] FOREIGN KEY ([Email_UUID])
		REFERENCES [Email] ([UUID])
;
GO


-----------------------------------------------------------
-- UserCredential Table --
-----------------------------------------------------------

ALTER TABLE [UserCredential]
	ADD
	CONSTRAINT [FK_UserCredential_UserUUID] FOREIGN KEY ([User_UUID])
		REFERENCES [User] ([UUID]),
	CONSTRAINT [FK_UserCredential_CredentialUUID] FOREIGN KEY ([Credential_UUID])
		REFERENCES [Credential] ([UUID])
;
GO




-------------------------------------------------------------------------------
-- Create Indexes --
-------------------------------------------------------------------------------

-----------------------------------------------------------
-- User Table --
-----------------------------------------------------------

-- Index for the Username column
CREATE INDEX [IX_User_Username]
	ON [User] ([Username])
;
GO

-----------------------------------------------------------
-- UserEmail Table --
-----------------------------------------------------------

-- Index for the User_UUID column
CREATE INDEX [IX_UserEmail_UserUUID]
	ON [UserEmail] ([User_UUID])
;
GO

-----------------------------------------------------------
-- Email Table --
-----------------------------------------------------------

-- Index for the Email column
CREATE INDEX [IX_Account_Email]
	ON [Email] ([Email])
;
GO

-----------------------------------------------------------
-- UserCredential Table --
-----------------------------------------------------------

-- Index for the User_UUID column
CREATE INDEX [IX_UserCredential_UserUUID]
	ON [UserCredential] ([User_UUID])
;
GO

