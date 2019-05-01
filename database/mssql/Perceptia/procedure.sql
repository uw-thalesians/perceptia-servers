/*
	Title: Perceptia Database Procedures
	Version: 0.7.1
	Schema Version: 0.5.0
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
	2019/03/02, Chris, Update to reflect field Uuid change, 0.4.1
	2019/03/02, Chris, Add add,get,delete for email, 0.5.0
	2019/04/28, Chris, Add delete for user, 0.6.0
	2019/04/28, Chris, Add Update Hash, DisplayName, FullName for user, 0.7.0
	2019/04/28, Chris, Fix params for Update sp, 0.7.1
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
-- GetUserEmailByUuid --
-----------------------------------------------------------

CREATE PROCEDURE [USP_GetUserEmailByUuid]
@Uuid UNIQUEIDENTIFIER
AS
SET NOCOUNT ON
;
BEGIN
	SELECT [E].[Email] FROM [User] AS [U]
		INNER JOIN [UserEmail] AS [UE]
		ON [U].[Uuid] = [UE].[User_UUID]
		INNER JOIN [Email] AS [E]
		ON [UE].[Email_Uuid] = [E].[Uuid]
		WHERE [U].[Uuid] = @Uuid
	;
END
;
GO

-----------------------------------------------------------
-- GetUserEmailByUsername --
-----------------------------------------------------------

CREATE PROCEDURE [USP_GetUserEmailByUsername]
@Username NVARCHAR(255)
AS
SET NOCOUNT ON
;
BEGIN
	SELECT [E].[Email] FROM [User] AS [U]
		INNER JOIN [UserEmail] AS [UE]
		ON [U].[Uuid] = [UE].[User_UUID]
		INNER JOIN [Email] AS [E]
		ON [UE].[Email_Uuid] = [E].[Uuid]
		WHERE [U].[Username] = @Username
	;
END
;
GO

-----------------------------------------------------------
-- GetUsernameByEmail --
-----------------------------------------------------------

CREATE PROCEDURE [USP_GetUsernameByEmail]
@Email NVARCHAR(255)
AS
SET NOCOUNT ON
;
BEGIN
	SELECT [U].[Username] FROM [User] AS [U]
		INNER JOIN [UserEmail] AS [UE]
		ON [U].[Uuid] = [UE].[User_UUID]
		INNER JOIN [Email] AS [E]
		ON [UE].[Email_Uuid] = [E].[Uuid]
		WHERE [E].[Email] = @Email
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
		SET @Uuid = NEWID()
	;
	DECLARE @CredentialUuid UNIQUEIDENTIFIER;
	BEGIN TRANSACTION [T1]
		INSERT INTO [User]
		([Uuid], [Username], [FullName], [DisplayName])
		VALUES
		(@Uuid, @Username, @FullName, @DisplayName)
		;	

		SET @CredentialUuid = NEWID()
		;

		INSERT INTO [Credential]
		([Uuid],[EncodedHash])
		VALUES
		(@CredentialUuid, @EncodedHash)
		;

		INSERT INTO [UserCredential]
		([User_Uuid], [Credential_Uuid])
		VALUES
		(@Uuid, @CredentialUuid)
		;
	COMMIT TRANSACTION [T1]
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


-----------------------------------------------------------
-- AddUserEmail --
-----------------------------------------------------------

CREATE PROCEDURE [USP_AddUserEmail]
@Uuid UNIQUEIDENTIFIER,
@Email NVARCHAR(255)
AS
SET NOCOUNT ON
;
BEGIN TRY
	IF @Uuid IS NULL
		THROW 50101, N'uuid must not be null', 1
		;
	IF @Email IS NULL
		THROW 50102, N'email must not be null', 1
		;
	IF NOT EXISTS (SELECT [Uuid] FROM [User] WHERE [Uuid] = @Uuid)
		THROW 50201, N'user does not exist', 1
		;
	;
	DECLARE @EmailUuid UNIQUEIDENTIFIER;
	BEGIN TRANSACTION [T1]
		SET @EmailUuid = NEWID();
		INSERT INTO [Email]
		([Uuid], [Email])
		VALUES
		(@EmailUuid, @Email)
		;

		INSERT INTO [UserEmail]
		([User_Uuid], [Email_Uuid])
		VALUES
		(@Uuid, @EmailUuid)
		;
	COMMIT TRANSACTION [T1]
	;
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

-----------------------------------------------------------
-- DeleteUserEmail --
-----------------------------------------------------------

CREATE PROCEDURE [USP_DeleteUserEmail]
@Uuid UNIQUEIDENTIFIER,
@Email NVARCHAR(255)
AS
SET NOCOUNT ON
;
BEGIN TRY
	IF @Uuid IS NULL
		THROW 50101, N'uuid must not be null', 1
		;
	IF @Email IS NULL
		THROW 50102, N'email must not be null', 1
		;
	IF NOT EXISTS (SELECT [Uuid] FROM [User] WHERE [Uuid] = @Uuid)
		THROW 50201, N'user does not exist', 1
		;
	;
	DECLARE @EmailUuid UNIQUEIDENTIFIER
	;
	SET @EmailUuid = (SELECT [E].[Uuid] FROM [User] AS [U]
					INNER JOIN [UserEmail] AS [UE]
					ON [U].[Uuid] = [UE].[User_UUID]
					INNER JOIN [Email] AS [E]
					ON [UE].[Email_Uuid] = [E].[Uuid]
					WHERE [U].[Uuid] = @Uuid AND [E].[Email] = @Email)
	;
	IF @EmailUuid IS NULL
		THROW 50202, N'user email does not exist', 1
		;
	;
	BEGIN TRANSACTION [T1]
		DELETE FROM [Email]
		WHERE [Uuid] = @EmailUuid
	COMMIT TRANSACTION [T1]
	;
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

-----------------------------------------------------------
-- DeleteUser --
-----------------------------------------------------------

CREATE PROCEDURE [USP_DeleteUser]
@Uuid UNIQUEIDENTIFIER
AS
SET NOCOUNT ON
;
BEGIN TRY
	IF @Uuid IS NULL
		THROW 50101, N'uuid must not be null', 1
		;
	IF NOT EXISTS (SELECT [Uuid] FROM [User] WHERE [Uuid] = @Uuid)
		THROW 50201, N'user does not exist', 1
		;
	BEGIN TRANSACTION [T1]
		DELETE FROM [User]
		WHERE [Uuid] = @Uuid
	COMMIT TRANSACTION [T1]
	;
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

-----------------------------------------------------------
-- UpdateUserEncodedHash --
-----------------------------------------------------------

CREATE PROCEDURE [USP_UpdateUserEncodedHash]
@Uuid UNIQUEIDENTIFIER
,@EncodedHash NVARCHAR(500)
AS
SET NOCOUNT ON
;
BEGIN TRY
	IF @Uuid IS NULL
		THROW 50101, N'uuid must not be null', 1
		;
	IF NOT EXISTS (SELECT [Uuid] FROM [User] WHERE [Uuid] = @Uuid)
		THROW 50201, N'user does not exist', 1
		;
	DECLARE @CredentialUuid UNIQUEIDENTIFIER
	;
	BEGIN TRANSACTION [T1]
		DELETE FROM [Credential]
		WHERE [Uuid] = (SELECT [C].[Uuid] FROM [Credential] AS [C]
			INNER JOIN [UserCredential] AS [UC]
				ON [C].[Uuid] = [UC].[Credential_Uuid]
			WHERE [UC].[User_Uuid] = @Uuid)
		;

		SET @CredentialUuid = NEWID()
		;

		INSERT INTO [Credential]
		([Uuid],[EncodedHash])
		VALUES
		(@CredentialUuid, @EncodedHash)
		;

		INSERT INTO [UserCredential]
		([User_Uuid], [Credential_Uuid])
		VALUES
		(@Uuid, @CredentialUuid)
		;
		
	COMMIT TRANSACTION [T1]
	;
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

-----------------------------------------------------------
-- UpdateUserFullName --
-----------------------------------------------------------

CREATE PROCEDURE [USP_UpdateUserFullName]
@Uuid UNIQUEIDENTIFIER
,@FullName NVARCHAR(255)
AS
SET NOCOUNT ON
;
BEGIN TRY
	IF @Uuid IS NULL
		THROW 50101, N'uuid must not be null', 1
		;
	IF NOT EXISTS (SELECT [Uuid] FROM [User] WHERE [Uuid] = @Uuid)
		THROW 50201, N'user does not exist', 1
		;
	BEGIN TRANSACTION [T1]
		UPDATE [User]
			SET [FullName] = @FullName
			WHERE [Uuid] = @Uuid
		;
	COMMIT TRANSACTION [T1]
	;
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

-----------------------------------------------------------
-- UpdateUserDisplayName --
-----------------------------------------------------------

CREATE PROCEDURE [USP_UpdateUserDisplayName]
@Uuid UNIQUEIDENTIFIER
,@DisplayName NVARCHAR(255)
AS
SET NOCOUNT ON
;
BEGIN TRY
	IF @Uuid IS NULL
		THROW 50101, N'uuid must not be null', 1
		;
	IF NOT EXISTS (SELECT [Uuid] FROM [User] WHERE [Uuid] = @Uuid)
		THROW 50201, N'user does not exist', 1
		;
	BEGIN TRANSACTION [T1]
		UPDATE [User]
			SET [DisplayName] = @DisplayName
			WHERE [Uuid] = @Uuid
		;
	COMMIT TRANSACTION [T1]
	;
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