/*
	Title: Perceptia Database Procedures
	Version: 0.3.0
*/
-------------------------------------------------------------------------------
-- Change Log --
-------------------------------------------------------------------------------
/*
	Date, Changer, Short Description, Version
	2019/02/19, Chris, Created Procedure, 0.1.0
	2019/02/19, Chris, Created InsertNewAccount, 0.2.0
	2019/02/19, Chris, Add Get procedures, 0.3.0
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
-- GetAccountInfoByUUID --
-----------------------------------------------------------

CREATE PROCEDURE [USP_GetAccountByUUID]
@AccountUUID UNIQUEIDENTIFIER
AS
SET NOCOUNT ON
;
BEGIN
	SELECT [AccountUUID], [Username], [FullName], [DisplayName], [Email]
	FROM [Account]
	WHERE [AccountUUID] = @AccountUUID
	;
END
;
GO

-----------------------------------------------------------
-- GetAccountInfoByUsername --
-----------------------------------------------------------

CREATE PROCEDURE [USP_GetAccountByUsername]
@Username NVARCHAR(255)
AS
SET NOCOUNT ON
;
BEGIN
	SELECT [AccountUUID], [Username], [FullName], [DisplayName], [Email]
	FROM [Account]
	WHERE [Username] = @Username
	;
END
;
GO

-----------------------------------------------------------
-- GetAccountInfoByEmail --
-----------------------------------------------------------

CREATE PROCEDURE [USP_GetAccountByEmail]
@Email NVARCHAR(255)
AS
SET NOCOUNT ON
;
BEGIN
	SELECT [AccountUUID], [Username], [FullName], [DisplayName], [Email]
	FROM [Account]
	WHERE [Email] = @Email
	;
END
;
GO

-----------------------------------------------------------
-- GetEncodedPasswordByUserName --
-----------------------------------------------------------

CREATE PROCEDURE [USP_GetEncodedPasswordByUserName]
@UserName NVARCHAR(255)
AS
SET NOCOUNT ON
;
BEGIN
	SELECT [EncodedPassword]
	FROM [Account]
	WHERE [UserName] = @UserName
	;
END
;
GO

-----------------------------------------------------------
-- InsertNewAccount --
-----------------------------------------------------------

CREATE PROCEDURE [USP_InsertNewAccount]
@AccountUUID UNIQUEIDENTIFIER,
@Username NVARCHAR(255),
@FullName NVARCHAR(255),
@DisplayName NVARCHAR(255),
@Email NVARCHAR(255),
@EncodedPassword NVARCHAR(500)
AS
SET NOCOUNT ON
;
BEGIN TRY
	BEGIN TRANSACTION [T1]
		INSERT INTO [Account]
		([AccountUUID], [Username], [FullName], [DisplayName], [Email], [EncodedPassword])
		VALUES
		(@AccountUUID, @Username, @FullName, @DisplayName, @Email, @EncodedPassword)
		;
			
	COMMIT TRANSACTION [T1]
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
