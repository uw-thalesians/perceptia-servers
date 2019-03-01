/*
	Title: Perceptia Database Procedures
	Version: 0.4.0
*/
-------------------------------------------------------------------------------
-- Change Log --
-------------------------------------------------------------------------------
/*
	Date, Changer, Short Description, Version
	2019/02/19, Chris, Created Procedure, 0.1.0
	2019/02/19, Chris, Created InsertNewAccount, 0.2.0
	2019/02/19, Chris, Add Get procedures, 0.3.0
	2019/02/28, Chris, Update procs for new schema, 0.4.0
*/

-------------------------------------------------------------------------------
-- Sections --
-------------------------------------------------------------------------------

-- Setup
-- Create Procedures

-------------------------------------------------------------------------------
-- Setup --
-------------------------------------------------------------------------------

USE [Perceptia]
;
GO

-------------------------------------------------------------------------------
-- Create Procedures --
-------------------------------------------------------------------------------

-----------------------------------------------------------
-- GetUserInfoByUUID --
-----------------------------------------------------------

CREATE PROCEDURE [USP_GetUserInfoByUUID]
@UUID UNIQUEIDENTIFIER
AS
SET NOCOUNT ON
;
BEGIN
	SELECT [UUID], [Username], [FullName], [DisplayName]
	FROM [User]
	WHERE [UUID] = @UUID
	;
END
;
GO

-----------------------------------------------------------
-- GetUserInfoByUsername --
-----------------------------------------------------------

CREATE PROCEDURE [USP_GetUserInfoByUsername]
@Username NVARCHAR(255)
AS
SET NOCOUNT ON
;
BEGIN
	SELECT [UUID], [Username], [FullName], [DisplayName]
	FROM [User]
	WHERE [Username] = @Username
	;
END
;
GO

-----------------------------------------------------------
-- GetUserEncodedHashByUsername --
-----------------------------------------------------------

CREATE PROCEDURE [USP_GetUserEncodedHashByUsername]
@Username NVARCHAR(255)
AS
SET NOCOUNT ON
;
BEGIN
	SELECT [EncodedHash]
	FROM [Credential] AS [C]
		INNER JOIN [dbo].[UserCredential] AS [UC]
		ON [C].UUID=[UC].[Credential_UUID]
		INNER JOIN [dbo].[User] AS [U]
		ON [UC].[User_UUID]=[U].[UUID]
	WHERE [U].[Username] = @Username
	;
END
;
GO

-----------------------------------------------------------
-- CreateUser --
-----------------------------------------------------------

CREATE PROCEDURE [USP_CreateUser]
@UUID UNIQUEIDENTIFIER = NULL,
@Username NVARCHAR(255),
@FullName NVARCHAR(255),
@DisplayName NVARCHAR(255),
@EncodedHash NVARCHAR(500)
AS
SET NOCOUNT ON
;
BEGIN TRY
	IF @UUID IS NULL
	BEGIN
		SET @UUID = NEWID();
	END
	;
	DECLARE @CredentialUUID UNIQUEIDENTIFIER;
	BEGIN TRANSACTION [T1]
		INSERT INTO [User]
		([UUID], [Username], [FullName], [DisplayName])
		VALUES
		(@UUID, @Username, @FullName, @DisplayName)
		;
			
	COMMIT TRANSACTION [T1]
	;
	BEGIN TRANSACTION [T2]
		SET @CredentialUUID = NEWID();

		INSERT INTO [Credential]
		([UUID],[EncodedHash])
		VALUES
		(@CredentialUUID, @EncodedHash)
		;
	
	COMMIT TRANSACTION [T2]
	;
	BEGIN TRANSACTION [T3]
		INSERT INTO [UserCredential]
		([User_UUID], [Credential_UUID])
		VALUES
		(@UUID, @CredentialUUID)
		;
	COMMIT TRANSACTION [T3]
	;
	-- Return the newly inserted user
	BEGIN
		EXECUTE USP_GetUserInfoByUUID @UUID
		;
	END
END TRY
BEGIN CATCH
	IF @@TRANCOUNT > 0
	BEGIN
		ROLLBACK
		;
	END
	;
	THROW
END CATCH;
;
GO
