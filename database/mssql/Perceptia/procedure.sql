/*
	Title: Perceptia Database Procedures
	Version: 0.4.1
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
	2019/03/01, Chris, Update to reflect field Uuid change, 0.4.1
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
-- GetUserInfoByUuid --
-----------------------------------------------------------

CREATE PROCEDURE [USP_GetUserInfoByUuid]
@Uuid UNIQUEIDENTIFIER
AS
SET NOCOUNT ON
;
BEGIN
	SELECT [Uuid], [Username], [FullName], [DisplayName]
	FROM [User]
	WHERE [Uuid] = @Uuid
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
	SELECT [Uuid], [Username], [FullName], [DisplayName]
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
		ON [C].Uuid=[UC].[Credential_Uuid]
		INNER JOIN [dbo].[User] AS [U]
		ON [UC].[User_Uuid]=[U].[Uuid]
	WHERE [U].[Username] = @Username
	;
END
;
GO

-----------------------------------------------------------
-- CreateUser --
-----------------------------------------------------------

CREATE PROCEDURE [USP_CreateUser]
@Uuid UNIQUEIDENTIFIER = NULL,
@Username NVARCHAR(255),
@FullName NVARCHAR(255),
@DisplayName NVARCHAR(255),
@EncodedHash NVARCHAR(500)
AS
SET NOCOUNT ON
;
BEGIN TRY
	IF @Uuid IS NULL
	BEGIN
		SET @Uuid = NEWID();
	END
	;
	DECLARE @CredentialUuid UNIQUEIDENTIFIER;
	BEGIN TRANSACTION [T1]
		INSERT INTO [User]
		([Uuid], [Username], [FullName], [DisplayName])
		VALUES
		(@Uuid, @Username, @FullName, @DisplayName)
		;
			
	COMMIT TRANSACTION [T1]
	;
	BEGIN TRANSACTION [T2]
		SET @CredentialUuid = NEWID();

		INSERT INTO [Credential]
		([Uuid],[EncodedHash])
		VALUES
		(@CredentialUuid, @EncodedHash)
		;
	
	COMMIT TRANSACTION [T2]
	;
	BEGIN TRANSACTION [T3]
		INSERT INTO [UserCredential]
		([User_Uuid], [Credential_Uuid])
		VALUES
		(@Uuid, @CredentialUuid)
		;
	COMMIT TRANSACTION [T3]
	;
	-- Return the newly inserted user
	BEGIN
		EXECUTE USP_GetUserInfoByUuid @Uuid
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
