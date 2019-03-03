/*
	Title: Perceptia Database Schema
	Version: 0.3.2
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
	2019/03/02, Chris, Change field UUID to Uuid, 0.3.1
	2019/03/02, Chris, On delete cascade, 0.3.2
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
	[Uuid] UNIQUEIDENTIFIER DEFAULT(NEWID()) NOT NULL,
	[Username] NVARCHAR(255) NOT NULL,
	[FullName] NVARCHAR(255),
	[DisplayName] NVARCHAR(255),
	[Created] DATETIME DEFAULT(GETDATE()),
	CONSTRAINT [PK_User_Uuid] PRIMARY KEY ([Uuid]),
	CONSTRAINT [UQ_User_Username] UNIQUE ([Username])
)
;
GO

-----------------------------------------------------------
-- Email Table --
-----------------------------------------------------------
-- Summary: Store email addresses

CREATE TABLE [Email] (
	[Uuid] UNIQUEIDENTIFIER DEFAULT(NEWID()) NOT NULL,
	[Email] NVARCHAR(255) NOT NULL,
	[Created] DATETIME DEFAULT(GETDATE()),
	CONSTRAINT [PK_Email_Uuid] PRIMARY KEY ([Uuid])
)
;
GO

-----------------------------------------------------------
-- UserEmail Table --
-----------------------------------------------------------
-- Summary: Associates an email with a user

CREATE TABLE [UserEmail] (
	[Uuid] UNIQUEIDENTIFIER DEFAULT(NEWID()) NOT NULL,
	[User_Uuid] UNIQUEIDENTIFIER NOT NULL,
	[Email_Uuid] UNIQUEIDENTIFIER NOT NULL,
	CONSTRAINT [PK_UserEmail_Uuid] PRIMARY KEY ([Uuid]),
)
;
GO

-----------------------------------------------------------
-- Credential Table --
-----------------------------------------------------------
-- Summary: Store login credentials

CREATE TABLE [Credential] (
	[Uuid] UNIQUEIDENTIFIER DEFAULT(NEWID()) NOT NULL,
	[Created] DATETIME DEFAULT(GETDATE()),
	[EncodedHash] NVARCHAR(500) NOT NULL,
	CONSTRAINT [PK_Credential_Uuid] PRIMARY KEY ([Uuid])
)
;
GO

-----------------------------------------------------------
-- UserCredential Table --
-----------------------------------------------------------
-- Summary: Associates a user with a login credential

CREATE TABLE [UserCredential] (
	[Uuid] UNIQUEIDENTIFIER DEFAULT(NEWID()) NOT NULL,
	[User_Uuid] UNIQUEIDENTIFIER NOT NULL,
	[Credential_Uuid] UNIQUEIDENTIFIER NOT NULL,
	CONSTRAINT [PK_UserCredential_Uuid] PRIMARY KEY ([Uuid]),
	CONSTRAINT [UQ_UserCredential_UserUuid] UNIQUE ([User_Uuid]),
	CONSTRAINT [UQ_UserCredential_CredentialUuid] UNIQUE ([Credential_Uuid])
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
	CONSTRAINT [FK_UserEmail_UserUuid] FOREIGN KEY ([User_Uuid])
		REFERENCES [User] ([Uuid])
		ON DELETE CASCADE,
	CONSTRAINT [FK_UserEmail_EmailUuid] FOREIGN KEY ([Email_Uuid])
		REFERENCES [Email] ([Uuid])
		ON DELETE CASCADE
;
GO


-----------------------------------------------------------
-- UserCredential Table --
-----------------------------------------------------------

ALTER TABLE [UserCredential]
	ADD
	CONSTRAINT [FK_UserCredential_UserUuid] FOREIGN KEY ([User_Uuid])
		REFERENCES [User] ([Uuid])
		ON DELETE CASCADE,
	CONSTRAINT [FK_UserCredential_CredentialUuid] FOREIGN KEY ([Credential_Uuid])
		REFERENCES [Credential] ([Uuid])
		ON DELETE CASCADE
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

-- Index for the User_Uuid column
CREATE INDEX [IX_UserEmail_UserUuid]
	ON [UserEmail] ([User_Uuid])
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

-- Index for the User_Uuid column
CREATE INDEX [IX_UserCredential_UserUuid]
	ON [UserCredential] ([User_Uuid])
;
GO

