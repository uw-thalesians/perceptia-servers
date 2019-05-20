/*
	Title: Perceptia Database Schema
	Version: 0.7.0
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
	2019/04/28, Chris, Update sp to 0.6.0, 0.4.0
	2019/04/28, Chris, Update sp to 0.7.0, 0.5.0
	2019/04/28, Chris, Update sp to 0.7.1, 0.5.0
	2019/05/18, Chris, Add Session Version Profile table, 0.7.0
*/

-------------------------------------------------------------------------------
-- Sections --
-------------------------------------------------------------------------------

-- Setup Database
-- Create Database Tables
-- Create Database Foreign Key Constraints
-- Create Database Indexes
-- Create Database Roles
-- Create Database Users
-- Populate Database

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
	CONTAINMENT = PARTIAL
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
-- Version Table --
-----------------------------------------------------------
-- Summary: Store information about the version of the schema and stored procedures

CREATE TABLE [Version] (
	[Uuid] UNIQUEIDENTIFIER DEFAULT(NEWID()) NOT NULL
	,[Name] NVARCHAR(255) NOT NULL
	,[Version] NVARCHAR(255)
	,[Description] NVARCHAR(255)
	,[Update] NVARCHAR(255)
	,[Created] DATETIME DEFAULT(GETDATE())
	,CONSTRAINT [PK_Version_Uuid] PRIMARY KEY ([Uuid])
	,CONSTRAINT [UQ_Version_Name] UNIQUE ([Name])
)
;
GO

-----------------------------------------------------------
-- User Table --
-----------------------------------------------------------
-- Summary: Store information about a specific user

CREATE TABLE [User] (
	[Uuid] UNIQUEIDENTIFIER DEFAULT(NEWID()) NOT NULL
	,[Username] NVARCHAR(255) NOT NULL
	,[FullName] NVARCHAR(255)
	,[DisplayName] NVARCHAR(255)
	,[Created] DATETIME DEFAULT(GETDATE())
	,CONSTRAINT [PK_User_Uuid] PRIMARY KEY ([Uuid])
	,CONSTRAINT [UQ_User_Username] UNIQUE ([Username])
)
;
GO

-----------------------------------------------------------
-- Email Table --
-----------------------------------------------------------
-- Summary: Store email addresses

CREATE TABLE [Email] (
	[Uuid] UNIQUEIDENTIFIER DEFAULT(NEWID()) NOT NULL
	,[Email] NVARCHAR(255) NOT NULL
	,[Created] DATETIME DEFAULT(GETDATE())
	,CONSTRAINT [PK_Email_Uuid] PRIMARY KEY ([Uuid])
)
;
GO

-----------------------------------------------------------
-- UserEmail Table --
-----------------------------------------------------------
-- Summary: Associates an email with a user

CREATE TABLE [UserEmail] (
	[Uuid] UNIQUEIDENTIFIER DEFAULT(NEWID()) NOT NULL
	,[User_Uuid] UNIQUEIDENTIFIER NOT NULL
	,[Email_Uuid] UNIQUEIDENTIFIER NOT NULL
	,CONSTRAINT [PK_UserEmail_Uuid] PRIMARY KEY ([Uuid])
	,CONSTRAINT [UQ_UserEmail_EmailUuid] UNIQUE ([Email_Uuid])
)
;
GO

-----------------------------------------------------------
-- Credential Table --
-----------------------------------------------------------
-- Summary: Store login credentials

CREATE TABLE [Credential] (
	[Uuid] UNIQUEIDENTIFIER DEFAULT(NEWID()) NOT NULL
	,[Created] DATETIME DEFAULT(GETDATE())
	,[EncodedHash] NVARCHAR(500) NOT NULL
	,CONSTRAINT [PK_Credential_Uuid] PRIMARY KEY ([Uuid])
)
;
GO

-----------------------------------------------------------
-- UserCredential Table --
-----------------------------------------------------------
-- Summary: Associates a user with a login credential

CREATE TABLE [UserCredential] (
	[Uuid] UNIQUEIDENTIFIER DEFAULT(NEWID()) NOT NULL
	,[User_Uuid] UNIQUEIDENTIFIER NOT NULL
	,[Credential_Uuid] UNIQUEIDENTIFIER NOT NULL
	,CONSTRAINT [PK_UserCredential_Uuid] PRIMARY KEY ([Uuid])
	,CONSTRAINT [UQ_UserCredential_UserUuid] UNIQUE ([User_Uuid])
	,CONSTRAINT [UQ_UserCredential_CredentialUuid] UNIQUE ([Credential_Uuid])
)
;
GO

-----------------------------------------------------------
-- Session Table --
-----------------------------------------------------------
-- Summary: Store information about user sessions

CREATE TABLE [Session] (
	[Uuid] UNIQUEIDENTIFIER DEFAULT(NEWID()) NOT NULL
	,[SessionId] NVARCHAR(255) NOT NULL
	,[Status] NVARCHAR(255) NOT NULL
	,[Created] DATETIME DEFAULT(GETDATE())
	,CONSTRAINT [PK_Session_Uuid] PRIMARY KEY ([Uuid])
)
;
GO

-----------------------------------------------------------
-- UserSession Table --
-----------------------------------------------------------
-- Summary: Associates a user with a session

CREATE TABLE [UserSession] (
	[Uuid] UNIQUEIDENTIFIER DEFAULT(NEWID()) NOT NULL
	,[User_Uuid] UNIQUEIDENTIFIER NOT NULL
	,[Session_Uuid] UNIQUEIDENTIFIER NOT NULL
	,CONSTRAINT [PK_UserSession_Uuid] PRIMARY KEY ([Uuid])
	,CONSTRAINT [UQ_UserSession_SessionUuid] UNIQUE ([Session_Uuid])
)
;
GO

-----------------------------------------------------------
-- Profile Table --
-----------------------------------------------------------
-- Summary: Store information for a user profile in the system

CREATE TABLE [Profile] (
	[Uuid] UNIQUEIDENTIFIER DEFAULT(NEWID()) NOT NULL
	,[Bio] NVARCHAR(1000)
	,[GravitarUrl] NVARCHAR(1000)
	,[Created] DATETIME DEFAULT(GETDATE())
	,CONSTRAINT [PK_Profile_Uuid] PRIMARY KEY ([Uuid])
)
;
GO

-----------------------------------------------------------
-- UserProfile Table --
-----------------------------------------------------------
-- Summary: Associates a profile with the user

CREATE TABLE [UserProfile] (
	[Uuid] UNIQUEIDENTIFIER DEFAULT(NEWID()) NOT NULL
	,[User_Uuid] UNIQUEIDENTIFIER NOT NULL
	,[Profile_Uuid] UNIQUEIDENTIFIER NOT NULL
	,CONSTRAINT [PK_UserProfile_Uuid] PRIMARY KEY ([Uuid])
	,CONSTRAINT [UQ_UserProfile_ProfileUuid] UNIQUE ([Profile_Uuid])
	,CONSTRAINT [UQ_UserProfile_UserUuid] UNIQUE ([User_Uuid])
)
;
GO

-----------------------------------------------------------
-- Profile Sharing Table --
-----------------------------------------------------------
-- Summary: Store information to indicate which fields can be shared

CREATE TABLE [ProfileSharing] (
	[Uuid] UNIQUEIDENTIFIER DEFAULT(NEWID()) NOT NULL
	,[Bio] NCHAR(1)
	,[GravitarUrl] NCHAR(1)
	,[DisplayName] NCHAR(1)
	,CONSTRAINT [PK_ProfileSharing_Uuid] PRIMARY KEY ([Uuid])
)
;
GO

-----------------------------------------------------------
-- UserProfileSharing Table --
-----------------------------------------------------------
-- Summary: Associates a profile sharing with the user

CREATE TABLE [UserProfileSharing] (
	[Uuid] UNIQUEIDENTIFIER DEFAULT(NEWID()) NOT NULL
	,[User_Uuid] UNIQUEIDENTIFIER NOT NULL
	,[ProfileSharing_Uuid] UNIQUEIDENTIFIER NOT NULL
	,CONSTRAINT [PK_UserProfileSharing_Uuid] PRIMARY KEY ([Uuid])
	,CONSTRAINT [UQ_UserProfileSharing_ProfileSharingUuid] UNIQUE ([ProfileSharing_Uuid])
	,CONSTRAINT [UQ_UserProfileSharing_UserUuid] UNIQUE ([User_Uuid])
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
		ON DELETE CASCADE
	,CONSTRAINT [FK_UserEmail_EmailUuid] FOREIGN KEY ([Email_Uuid])
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
		ON DELETE CASCADE
	,CONSTRAINT [FK_UserCredential_CredentialUuid] FOREIGN KEY ([Credential_Uuid])
		REFERENCES [Credential] ([Uuid])
		ON DELETE CASCADE
;
GO

-----------------------------------------------------------
-- UserSession Table --
-----------------------------------------------------------

ALTER TABLE [UserSession]
	ADD
	CONSTRAINT [FK_UserSession_UserUuid] FOREIGN KEY ([User_Uuid])
		REFERENCES [User] ([Uuid])
		ON DELETE CASCADE
	,CONSTRAINT [FK_UserSession_SessionUuid] FOREIGN KEY ([Session_Uuid])
		REFERENCES [Session] ([Uuid])
		ON DELETE CASCADE
;
GO

-----------------------------------------------------------
-- UserProfile Table --
-----------------------------------------------------------

ALTER TABLE [UserProfile]
	ADD
	CONSTRAINT [FK_UserProfile_UserUuid] FOREIGN KEY ([User_Uuid])
		REFERENCES [User] ([Uuid])
		ON DELETE CASCADE
	,CONSTRAINT [FK_UserProfile_ProfileUuid] FOREIGN KEY ([Profile_Uuid])
		REFERENCES [Profile] ([Uuid])
		ON DELETE CASCADE
;
GO

-----------------------------------------------------------
-- UserProfileSharing Table --
-----------------------------------------------------------

ALTER TABLE [UserProfileSharing]
	ADD
	CONSTRAINT [FK_UserProfileSharing_UserUuid] FOREIGN KEY ([User_Uuid])
		REFERENCES [User] ([Uuid])
		ON DELETE CASCADE
	,CONSTRAINT [FK_UserProfileSharing_ProfileSharingUuid] FOREIGN KEY ([ProfileSharing_Uuid])
		REFERENCES [ProfileSharing] ([Uuid])
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

-------------------------------------------------------------------------------
-- Create Roles --
-------------------------------------------------------------------------------

-----------------------------------------------------------
-- Execute Stored Procedures in DB--
-----------------------------------------------------------
CREATE ROLE [RL_ExecuteAllProcedures]
;
GO

GRANT EXECUTE TO [RL_ExecuteAllProcedures]
;
GO

-------------------------------------------------------------------------------
-- Populate Database --
-------------------------------------------------------------------------------

-----------------------------------------------------------
-- Version Table --
-----------------------------------------------------------

INSERT INTO [Version] (
		[Uuid]
		,[Name]
		,[Version]
		,[Description]
	)
VALUES (
		N'8FBE90DA-70C2-4C0C-91AB-A2B8FE31F0D4'
		,N'Schema'
		,N'0.6.0'
		,N'The Perceptia Database Schema.'
	)
;
GO