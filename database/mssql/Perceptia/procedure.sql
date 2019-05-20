/*
	Title: Perceptia Database Procedures
	Version: 0.8.0
	Schema Version: 0.7.0
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
	2019/05/18, Chris, Add Session get and delete sp, 0.8.0
*/

-------------------------------------------------------------------------------
-- Sections --
-------------------------------------------------------------------------------

-- Setup
-- Populate Database
-- Procedure Error Notes
-- Create Procedures

-------------------------------------------------------------------------------
-- Setup --
-------------------------------------------------------------------------------

USE [Perceptia]
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
		N'1F51BCCE-959B-4732-97D4-3AD688850ED8'
		,N'Stored Procedures'
		,N'0.8.0'
		,N'The Perceptia Database Stored Procedures.'
	)
;
GO

-------------------------------------------------------------------------------
-- Procedure Error Notes --
-------------------------------------------------------------------------------

/*
	Range: 50000-50999
	Meaning:
		50100s: Required Value not provided (value was null)
		50200s: Provided Value not valid (invalid syntax/format)
		50300s: Referenced Object does not exist (identifier provided didn't
				match existing object, object not found)
		50400s: Referenced Object already exists (identifier provided was already
				found in the system, conflict with existing object)

	Meaning of specific value within range depends on procedure

*/

-------------------------------------------------------------------------------
-- Create Procedures --
-------------------------------------------------------------------------------

----------------------------------------------------------------
-------- READ Procedures --------
----------------------------------------------------------------

-----------------------------------------------------------
-- ReadProcedureVersion --
-----------------------------------------------------------

-- USP_ReadProcedureVersion gets the version of the Stored Procedures applied.
-- Parameters none
-- Outputs
--	Query row containing 1 columns (should be exactly one row).
--		Version: NVARCHAR(255) version of stored procedures available.
-- Errors none
CREATE PROCEDURE [USP_ReadProcedureVersion]
AS
SET NOCOUNT ON
;
BEGIN
	SELECT [Version]
		FROM [Version]
		WHERE [Name] = N'Stored Procedures'
	;
END
;
GO

-----------------------------------------------------------
-- ReadUserInfoByUuid --
-----------------------------------------------------------

-- USP_ReadUserInfoByUuid gets basic informaiton about the given user.
-- Parameters
--	@UserUuid:	UNIQUEIDENTIFIER the uuid for the user who's information should be returned.
--				Must be a valid v4 UUID.
-- Outputs
--	Query row containing 3 columns (should be exactly one row).
--		Uuid: UNIQUEIDENTIFIER the uuid of the user who's info was requested.
--		Username: NVARCHAR(255) the username of the user who's info was requested.
--		DisplayName: NVARCHAR(255) the DisplayName of the user who's info was requested.
-- Errors
--	50101: The provided UserUuid was null.
--	50301: No user found with the provided UserUuid.
CREATE PROCEDURE [USP_ReadUserInfoByUuid]
	@UserUuid UNIQUEIDENTIFIER
AS
SET NOCOUNT ON
;
BEGIN
	IF @UserUuid IS NULL
		THROW 50101, N'uuid must not be null', 1
	;
	IF NOT EXISTS (SELECT [Uuid] FROM [User] WHERE [Uuid] = @UserUuid)
		THROW 50301, N'user does not exist', 1
	;

	SELECT [Uuid], [Username], [DisplayName]
		FROM [User]
		WHERE [Uuid] = @UserUuid
	;
END
;
GO

-----------------------------------------------------------
-- ReadUserProfileByUuid --
-----------------------------------------------------------

-- USP_ReadUserProfileByUuid gets profile informaiton for a given user.
-- Parameters
--	@UserUuid:	UNIQUEIDENTIFIER the uuid for the user who's information should be returned.
--				Must be a valid v4 UUID.
-- Outputs
--	Query row containing 5 columns (should be exactly one row).
--		UserUuid: UNIQUEIDENTIFIER the uuid of the user who's info was requested.
--		Username: NVARCHAR(255) the username of the user who's info was requested.
--		DisplayName: NVARCHAR(255) the DisplayName of the user who's info was requested.
--		Bio: NVARCHAR(1000) a user provided description of self.
--		GravitarUrl: NVARCHAR(1000) avitar image for user.
-- Errors
--	50101: The provided UserUuid was null.
--	50301: No user found with the provided UserUuid.
CREATE PROCEDURE [USP_ReadUserProfileByUuid]
	@UserUuid UNIQUEIDENTIFIER
AS
SET NOCOUNT ON
;
BEGIN
	IF @UserUuid IS NULL
		THROW 50101, N'uuid must not be null', 1
	;
	IF NOT EXISTS (SELECT [Uuid] FROM [User] WHERE [Uuid] = @UserUuid)
		THROW 50301, N'user does not exist', 1
	;

	SELECT [U].[Uuid] AS [UserUuid], [U].[Username], [U].[DisplayName], [P].[Bio], [P].[GravitarUrl], [PS].[DisplayName] AS [ShareDisplayName], [PS].[Bio] AS [ShareBio], [PS].[GravitarUrl] AS [ShareGravitarUrl]
		FROM [User] AS [U]
		INNER JOIN [UserProfile] AS [UP]
			ON [U].[Uuid] = [UP].[User_Uuid]
		INNER JOIN [Profile] AS [P]
			ON [UP].[Profile_Uuid] = [P].[Uuid]
		INNER JOIN [UserProfileSharing] AS [UPS]
			ON [U].[Uuid] = [UPS].[User_Uuid]
		INNER JOIN [ProfileSharing] AS [PS]
			ON [UPS].[ProfileSharing_Uuid] = [PS].[Uuid]
		WHERE [U].[Uuid] = @UserUuid
	;
END
;
GO

-----------------------------------------------------------
-- ReadUserUuidByUsername --
-----------------------------------------------------------

-- USP_ReadUserUuidByUsername gets the uuid for a user based on the provided username.
-- Parameters
--	@Username:	NVARCHAR(255) the username for the user who's uuid should be returned.
--				Must be a valid username in the system.
-- Outputs
--	Query row containing 3 columns (should be exactly one row)
--		Uuid: UNIQUEIDENTIFIER the uuid of the user who's info was requested.
-- Errors
--	50101: The provided Username was null
--	50301: No user found with the provided Username
CREATE PROCEDURE [USP_ReadUserUuidByUsername]
	@Username NVARCHAR(255)
AS
SET NOCOUNT ON
;
BEGIN
	IF @Username IS NULL
		THROW 50101, N'username must not be null', 1
	;
	IF NOT EXISTS (SELECT [Uuid] FROM [User] WHERE [Username] = @Username)
		THROW 50301, N'user does not exist', 1
	;
	SELECT [Uuid]
		FROM [User]
		WHERE [Username] = @Username
	;
END
;
GO

-----------------------------------------------------------
-- ReadUserDisplayNameByUuid --
-----------------------------------------------------------

-- USP_ReadUserProfileByUuid gets the displayname for the given user.
-- Parameters
--	@UserUuid:	UNIQUEIDENTIFIER the uuid for the user who's information should be returned.
--				Must be a valid v4 UUID.
-- Outputs
--	Query row containing 1 columns (should be exactly one row).
--		DisplayName: NVARCHAR(255) the DisplayName of the user who's info was requested.
-- Errors
--	50101: The provided UserUuid was null.
--	50301: No user found with the provided UserUuid.
CREATE PROCEDURE [USP_ReadUserDisplayNameByUuid]
	@UserUuid UNIQUEIDENTIFIER
AS
SET NOCOUNT ON
;
BEGIN
	IF @UserUuid IS NULL
		THROW 50101, N'uuid must not be null', 1
	;
	IF NOT EXISTS (SELECT [Uuid] FROM [User] WHERE [Uuid] = @UserUuid)
		THROW 50301, N'user does not exist', 1
	;

	SELECT [DisplayName]
		FROM [User]
		WHERE [Uuid] = @UserUuid
	;
END
;
GO

-----------------------------------------------------------
-- ReadUserFullNameByUuid --
-----------------------------------------------------------

-- USP_ReadUserFullNameByUuid gets the FullName for the given user.
-- Parameters
--	@UserUuid:	UNIQUEIDENTIFIER the uuid for the user who's information should be returned.
--				Must be a valid v4 UUID.
-- Outputs
--	Query row containing 1 columns (should be exactly one row).
--		FullName: NVARCHAR(255) the FullName of the user who's info was requested.
-- Errors
--	50101: The provided UserUuid was null.
--	50301: No user found with the provided UserUuid.
CREATE PROCEDURE [USP_ReadUserFullNameByUuid]
	@UserUuid UNIQUEIDENTIFIER
AS
SET NOCOUNT ON
;
BEGIN
	IF @UserUuid IS NULL
		THROW 50101, N'uuid must not be null', 1
	;
	IF NOT EXISTS (SELECT [Uuid] FROM [User] WHERE [Uuid] = @UserUuid)
		THROW 50301, N'user does not exist', 1
	;

	SELECT [FullName]
		FROM [User]
		WHERE [Uuid] = @UserUuid
	;
END
;
GO

-----------------------------------------------------------
-- ReadUserUsernameByUuid --
-----------------------------------------------------------

-- USP_ReadUserUsernameByUuid gets the username for the given user.
-- Parameters
--	@UserUuid:	UNIQUEIDENTIFIER the uuid for the user who's information should be returned.
--				Must be a valid v4 UUID.
-- Outputs
--	Query row containing 1 columns (should be exactly one row).
--		Username: NVARCHAR(255) the Username of the user who's info was requested.
-- Errors
--	50101: The provided UserUuid was null.
--	50301: No user found with the provided UserUuid.
CREATE PROCEDURE [USP_ReadUserUsernameByUuid]
	@UserUuid UNIQUEIDENTIFIER
AS
SET NOCOUNT ON
;
BEGIN
	IF @UserUuid IS NULL
		THROW 50101, N'uuid must not be null', 1
	;
	IF NOT EXISTS (SELECT [Uuid] FROM [User] WHERE [Uuid] = @UserUuid)
		THROW 50301, N'user does not exist', 1
	;

	SELECT [Username]
		FROM [User]
		WHERE [Uuid] = @UserUuid
	;
END
;
GO

-----------------------------------------------------------
-- ReadUserEmailsByUuid --
-----------------------------------------------------------

-- USP_ReadUserEmailsByUuid returns a list of the emails associated with the given user.
-- Parameters
--	@UserUuid:	UNIQUEIDENTIFIER the UserUuid for the user who's information should be returned.
--				Must be a valid v4 UUID.
-- Outputs
--	Query row containing 3 columns (may return 0 or more rows).
--		Uuid: UNIQUEIDENTIFIER of Email.
--		Email: NVARCHAR(255) a single email.
--		Created: DATETIME the date when the email was added.
-- Errors
--	50101: The provided UserUuid was null.
--	50301: No user found with the provided UserUuid.
CREATE PROCEDURE [USP_ReadUserEmailsByUuid]
	@UserUuid UNIQUEIDENTIFIER
AS
SET NOCOUNT ON
;
BEGIN
	IF @UserUuid IS NULL
		THROW 50101, N'uuid must not be null', 1
	;
	IF NOT EXISTS (SELECT [Uuid] FROM [User] WHERE [Uuid] = @UserUuid)
		THROW 50301, N'user does not exist', 1
	;
	SELECT [E].[Uuid], [E].[Email], [E].[Created] 
		FROM [User] AS [U]
		INNER JOIN [UserEmail] AS [UE]
			ON [U].[Uuid] = [UE].[User_UUID]
		INNER JOIN [Email] AS [E]
			ON [UE].[Email_Uuid] = [E].[Uuid]
		WHERE [U].[Uuid] = @UserUuid
	;
END
;
GO

-----------------------------------------------------------
-- ReadUserEmailsByUsername --
-----------------------------------------------------------

-- USP_ReadUserEmailsByUsername returns a list of the emails associated with the given user.
-- Parameters
--	@Username:	the username for the user who's information should be returned.
--				Must be a valid username in the system.
-- Outputs
--	Query row containing 3 columns (may return 0 or more rows).
--		Uuid: UNIQUEIDENTIFIER of Email.
--		Email: NVARCHAR(255) a single email.
--		Created: DATETIME the date when the email was added.
-- Errors
--	50101: The provided Username was null.
--	50301: No user found with the provided Username.
CREATE PROCEDURE [USP_ReadUserEmailsByUsername]
	@Username NVARCHAR(255)
AS
SET NOCOUNT ON
;
BEGIN
	IF @Username IS NULL
		THROW 50101, N'username must not be null', 1
	;
	IF NOT EXISTS (SELECT [Uuid] FROM [User] WHERE [Username] = @Username)
		THROW 50301, N'user does not exist', 1
	;
	SELECT [E].[Uuid], [E].[Email], [E].[Created] 
		FROM [User] AS [U]
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
-- ReadUserUsernamesByEmail --
-----------------------------------------------------------

-- USP_ReadUserUsernamesByEmail returns a list of the usernames associated with the given email.
-- Parameters
--	@Email:	the username for the user who's information should be returned.
--				Must be a valid username in the system.
-- Outputs
--	Query row containing 3 columns (may return 0 or more rows).
--		Uuid: UNIQUEIDENTIFIER of a user.
--		Username: NVARCHAR(255) of a user.
--		Created: DATETIME the date when the user was created.
-- Errors
--	50101: The provided Email was null.
--	50301: Provided email not found in system.
CREATE PROCEDURE [USP_ReadUserUsernamesByEmail]
	@Email NVARCHAR(255)
AS
SET NOCOUNT ON
;
BEGIN
	IF @Email IS NULL
		THROW 50101, N'email must not be null', 1
	;
	IF NOT EXISTS (SELECT [Email] FROM [Email] WHERE [Email] = @Email)
		THROW 50301, N'user does not exist', 1
	;
	SELECT [U].[Uuid], [U].[Username], [U].[Created] 
		FROM [User] AS [U]
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
-- ReadUserEncodedHashByUsername --
-----------------------------------------------------------

-- USP_ReadUserEncodedHashByUsername returns the encoded has stored for the given user.
-- Parameters
--	@Username:	the username for the user who's encoded hash should be returned.
--				Must be a valid username in the system.
-- Outputs
--	Query row containing 1 column (should return exactly one row).
--		EncodedHash: NVARCHAR(500) is the encoded has for the user.
-- Errors
--	50101: The provided Username was null.
--	50301: No user found with the provided Username.
CREATE PROCEDURE [USP_ReadUserEncodedHashByUsername]
	@Username NVARCHAR(255)
AS
SET NOCOUNT ON
;
BEGIN
	IF @Username IS NULL
		THROW 50101, N'username must not be null', 1
	;
	IF NOT EXISTS (SELECT [Uuid] FROM [User] WHERE [Username] = @Username)
		THROW 50301, N'user does not exist', 1
	;
	SELECT [EncodedHash]
		FROM [dbo].[Credential] AS [C]
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
-- ReadUserSessionsByUuid --
-----------------------------------------------------------

-- USP_ReadUserSessionsByUuid returns a list of the sessions the user started.
-- Parameters
--	@UserUuid:	UNIQUEIDENTIFIER the UserUuid for the user who's information should be returned.
--				Must be a valid v4 UUID.
-- Outputs
--	Query row containing 4 columns (may return 0 or more rows).
--		Uuid: UNIQUEIDENTIFIER of session.
--		SessionId: NVARCHAR(255) session id portion of access token.
--		Status: NVARCHAR(255) is session active, one of {Active, Expired}.
--		Created: DATETIME the date when the session was added.
-- Errors
--	50101: The provided UserUuid was null.
--	50301: No user found with the provided UserUuid.
CREATE PROCEDURE [USP_ReadUserSessionsByUuid]
	@UserUuid UNIQUEIDENTIFIER
AS
SET NOCOUNT ON
;
BEGIN
	IF @UserUuid IS NULL
		THROW 50101, N'uuid must not be null', 1
	;
	IF NOT EXISTS (SELECT [Uuid] FROM [User] WHERE [Uuid] = @UserUuid)
		THROW 50301, N'user does not exist', 1
	;
	SELECT [S].[Uuid], [S].[SessionId], [S].[Status], [S].[Created]
		FROM [dbo].[Session] AS [S]
		INNER JOIN [dbo].[UserSession] AS [US]
			ON [S].Uuid=[US].[Session_Uuid]
		INNER JOIN [dbo].[User] AS [U]
			ON [US].[User_Uuid]=[U].[Uuid]
		WHERE [U].[Uuid] = @UserUuid
	;
END
;
GO

-----------------------------------------------------------
-- ReadUserActiveSessionsByUuid --
-----------------------------------------------------------

-- USP_ReadUserActiveSessionsByUuid returns a list of the sessions the user started
-- that are still active.
-- Parameters
--	@UserUuid:	UNIQUEIDENTIFIER the UserUuid for the user who's information should be returned.
--				Must be a valid v4 UUID.
-- Outputs
--	Query row containing 4 columns (may return 0 or more rows).
--		Uuid: UNIQUEIDENTIFIER of session.
--		SessionId: NVARCHAR(255) session id portion of access token.
--		Status: NVARCHAR(255) is session active, one of {Active, Expired}.
--		Created: DATETIME the date when the session was added.
-- Errors
--	50101: The provided UserUuid was null.
--	50301: No user found with the provided UserUuid.
CREATE PROCEDURE [USP_ReadUserActiveSessionsByUuid]
	@UserUuid UNIQUEIDENTIFIER
AS
SET NOCOUNT ON
;
BEGIN
	IF @UserUuid IS NULL
		THROW 50101, N'uuid must not be null', 1
	;
	IF NOT EXISTS (SELECT [Uuid] FROM [User] WHERE [Uuid] = @UserUuid)
		THROW 50301, N'user does not exist', 1
	;
	SELECT [S].[Uuid], [S].[SessionId], [S].[Status], [S].[Created]
		FROM [dbo].[Session] AS [S]
		INNER JOIN [dbo].[UserSession] AS [US]
			ON [S].Uuid=[US].[Session_Uuid]
		INNER JOIN [dbo].[User] AS [U]
			ON [US].[User_Uuid]=[U].[Uuid]
		WHERE [U].[Uuid] = @UserUuid AND [S].[Status] = N'Active'
	;
END
;
GO

----------------------------------------------------------------
-------- CREATE Procedures --------
----------------------------------------------------------------

-----------------------------------------------------------
-- CreateUser --
-----------------------------------------------------------

-- USP_CreateUser inserts the provided information, adding the user to the database.
-- Parameters
--	@UserUuid:	UNIQUEIDENTIFIER (optional) the UserUuid that the user should be created with.
--				Must be a valid v4 UUID. 
--	@Username:	NVARCHAR(255) (required) the username for the user who's should be added to the database.
--				Must be a valid username in the system.
--	@FullName:	NVARCHAR(255) (required) the FullName for the user who's should be added to the database.
--	@DisplayName:	NVARCHAR(255) (required) the DisplayName for the user who's should be added to the database.
--	@EncodedHash:	NVARCHAR(255) (required) the EncodedHash of the users password.
-- Outputs
--	Query row containing 4 columns (may return 0 or more rows).
--		Uuid: UNIQUEIDENTIFIER uuid of user added.
--		Username: NVARCHAR(255) username of user added.
--		FullName: NVARCHAR(255) full name of user added.
--		DisplayName: NVARCHAR(255) the display name of the user added.
-- Errors
--	50401: User with provided uuid already exists.
CREATE PROCEDURE [USP_CreateUser]
	@UserUuid UNIQUEIDENTIFIER = NULL
	,@Username NVARCHAR(255)
	,@FullName NVARCHAR(255)
	,@DisplayName NVARCHAR(255)
	,@EncodedHash NVARCHAR(500)
AS
SET NOCOUNT ON
;
BEGIN TRY
	IF @UserUuid IS NULL
		SET @UserUuid = NEWID()
	;
	IF EXISTS (SELECT [Uuid] FROM [User] WHERE [Uuid] = @UserUuid)
		THROW 50401, N'user with provided uuid already exists', 1
	;
	DECLARE @CredentialUuid UNIQUEIDENTIFIER;
	DECLARE @ProfileUuid UNIQUEIDENTIFIER;
	DECLARE @ProfileSharingUuid UNIQUEIDENTIFIER;
	BEGIN TRANSACTION [T1]
		INSERT INTO [User]
			([Uuid], [Username], [FullName], [DisplayName])
		VALUES
			(@UserUuid, @Username, @FullName, @DisplayName)
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
			(@UserUuid, @CredentialUuid)
		;
		SET @ProfileUuid = NEWID()
		INSERT INTO [Profile]
			([Uuid])
		VALUES
			(@ProfileUuid)
		;
		INSERT INTO [UserProfile]
			([User_Uuid], [Profile_Uuid])
		VALUES
			(@UserUuid, @ProfileUuid)
		;
		SET @ProfileSharingUuid = NEWID()
		INSERT INTO [ProfileSharing]
			([Uuid], [Bio], [GravitarUrl], [DisplayName])
		VALUES
			(@ProfileSharingUuid, N'N', N'N', N'N')
		;
		INSERT INTO [UserProfileSharing]
			([User_Uuid], [ProfileSharing_Uuid])
		VALUES
			(@UserUuid, @ProfileSharingUuid)
		;
	COMMIT TRANSACTION [T1]
	;
	-- Return the newly inserted user
	BEGIN
		SELECT [Uuid], [Username], [FullName], [DisplayName]
			FROM [User]
			WHERE [Uuid] = @UserUuid
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
END CATCH
;
GO


-----------------------------------------------------------
-- CreateUserEmail --
-----------------------------------------------------------

-- USP_CreateUserEmail adds the provided email to the specified users list of emails.
-- Parameters
--	@UserUuid:	UNIQUEIDENTIFIER the UserUuid for the user who's information should be inserted.
--				Must be a valid v4 UUID.
--	@Email: NVARCHAR(255) the email that should be added to the users account.
-- Outputs
--	Query result indicating the number of rows updated (should be 1)
-- Errors
--	50101: The provided UserUuid was null.
--	50102: The provided Email was null.
--	50301: No user found with the provided UserUuid.
--	50401: Provided email already in users list of emails.
CREATE PROCEDURE [USP_CreateUserEmail]
	@UserUuid UNIQUEIDENTIFIER
	,@Email NVARCHAR(255)
AS
SET NOCOUNT ON
;
BEGIN TRY
	IF @UserUuid IS NULL
		THROW 50101, N'uuid must not be null', 1
		;
	IF @Email IS NULL
		THROW 50102, N'email must not be null', 1
		;
	IF NOT EXISTS (SELECT [Uuid] FROM [User] WHERE [Uuid] = @UserUuid)
		THROW 50301, N'user does not exist', 1
		;
	IF EXISTS (SELECT [E].[Email] FROM [User]  AS [U]
		INNER JOIN [dbo].[UserEmail] AS [UE]
			ON [U].[Uuid] = [UE].[User_Uuid]
		INNER JOIN [dbo].[Email] AS [E]
			ON [UE].[Email_Uuid] = [E].[Uuid]
	WHERE [U].[Uuid] = @UserUuid AND [E].[Email] = @Email) 
		THROW 50401, N'email already exists for user', 1
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
			(@UserUuid, @EmailUuid)
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

----------------------------------------------------------------
-------- UPDATE Procedures --------
----------------------------------------------------------------

-----------------------------------------------------------
-- UpdateUserEncodedHash --
-----------------------------------------------------------

-- USP_UpdateUserEncodedHash replaces the existing EncodedHash with the provided one.
-- Parameters
--	@UserUuid:	UNIQUEIDENTIFIER the UserUuid for the user who's information should be inserted.
--				Must be a valid v4 UUID.
--	@EncodedHash NVARCHAR(500) the EncodedHash to be added.
-- Outputs
--	Query result indicating the number of rows updated (should be 1)
-- Errors
--	50101: The provided UserUuid was null.
--	50102: The provided EncodedHash was null.
--	50301: No user found with the provided UserUuid.
CREATE PROCEDURE [USP_UpdateUserEncodedHash]
	@UserUuid UNIQUEIDENTIFIER
	,@EncodedHash NVARCHAR(500)
AS
SET NOCOUNT ON
;
BEGIN TRY
	IF @UserUuid IS NULL
		THROW 50101, N'uuid must not be null', 1
		;
	IF @EncodedHash IS NULL
		THROW 50101, N'encoded hash must not be null', 1
		;
	IF NOT EXISTS (SELECT [Uuid] FROM [User] WHERE [Uuid] = @UserUuid)
		THROW 50301, N'user does not exist', 1
		;
	DECLARE @CredentialUuid UNIQUEIDENTIFIER
	;
	BEGIN TRANSACTION [T1]
		DELETE FROM [Credential]
			WHERE [Uuid] = (
					SELECT [C].[Uuid] FROM [Credential] AS [C]
						INNER JOIN [UserCredential] AS [UC]
							ON [C].[Uuid] = [UC].[Credential_Uuid]
						WHERE [UC].[User_Uuid] = @UserUuid
				)
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
			(@UserUuid, @CredentialUuid)
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

-- USP_UpdateUserFullName replaces the existing FullName with the provided one.
-- Parameters
--	@UserUuid:	UNIQUEIDENTIFIER the UserUuid for the user who's information should be updated.
--				Must be a valid v4 UUID.
--	@FullName NVARCHAR(255) the FullName to be added.
-- Outputs
--	Query result indicating the number of rows updated (should be 1)
-- Errors
--	50101: The provided UserUuid was null.
--	50102: The provided FullName was null.
--	50301: No user found with the provided UserUuid.
CREATE PROCEDURE [USP_UpdateUserFullName]
	@UserUuid UNIQUEIDENTIFIER
	,@FullName NVARCHAR(255)
AS
SET NOCOUNT ON
;
BEGIN TRY
	IF @UserUuid IS NULL
		THROW 50101, N'uuid must not be null', 1
		;
	IF @FullName IS NULL
		THROW 50102, N'full name must not be null', 1
		;
	IF NOT EXISTS (SELECT [Uuid] FROM [User] WHERE [Uuid] = @UserUuid)
		THROW 50301, N'user does not exist', 1
		;
	BEGIN TRANSACTION [T1]
		UPDATE [User]
			SET [FullName] = @FullName
			WHERE [Uuid] = @UserUuid
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
-- UpdateUserUsername --
-----------------------------------------------------------

-- USP_UpdateUserUsername replaces the existing username with the provided one.
-- Parameters
--	@UserUuid:	UNIQUEIDENTIFIER the UserUuid for the user who's information should be updated.
--				Must be a valid v4 UUID.
--	@Username NVARCHAR(255) the username to be added.
-- Outputs
--	Query result indicating the number of rows updated (should be 1)
-- Errors
--	50101: The provided UserUuid was null.
--	50102: The provided Username was null.
--	50301: No user found with the provided UserUuid.
--	50401: Provided username already in use.
--	50402: Provided username is users current username.
CREATE PROCEDURE [USP_UpdateUserUsername]
	@UserUuid UNIQUEIDENTIFIER
	,@Username NVARCHAR(255)
AS
SET NOCOUNT ON
;
BEGIN TRY
	IF @UserUuid IS NULL
		THROW 50101, N'uuid must not be null', 1
		;
	IF @Username IS NULL
		THROW 50102, N'full name must not be null', 1
		;
	IF NOT EXISTS (SELECT [Uuid] FROM [User] WHERE [Uuid] = @UserUuid)
		THROW 50301, N'user does not exist', 1
		;
	IF @Username = (SELECT [Username] FROM [User] WHERE [Uuid] = @UserUuid)
		THROW 50402, N'username already users username', 1
		;
	IF EXISTS (SELECT [Username] FROM [User] WHERE [Username] = @Username)
		THROW 50401, N'username already in use', 1
		;
	BEGIN TRANSACTION [T1]
		UPDATE [User]
			SET [Username] = @Username
			WHERE [Uuid] = @UserUuid
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

-- USP_UpdateUserDisplayName replaces the existing DisplayName with the provided one.
-- Parameters
--	@UserUuid:	UNIQUEIDENTIFIER the UserUuid for the user who's information should be updated.
--				Must be a valid v4 UUID.
--	@DisplayName NVARCHAR(255) the DisplayName to be added.
-- Outputs
--	Query result indicating the number of rows updated (should be 1)
-- Errors
--	50101: The provided UserUuid was null.
--	50102: The provided DisplayName was null.
--	50301: No user found with the provided UserUuid.
CREATE PROCEDURE [USP_UpdateUserDisplayName]
	@UserUuid UNIQUEIDENTIFIER
	,@DisplayName NVARCHAR(255)
AS
SET NOCOUNT ON
;
BEGIN TRY
	IF @UserUuid IS NULL
		THROW 50101, N'uuid must not be null', 1
		;
	IF @DisplayName IS NULL
		THROW 50102, N'display name must not be null', 1
		;
	IF NOT EXISTS (SELECT [Uuid] FROM [User] WHERE [Uuid] = @UserUuid)
		THROW 50301, N'user does not exist', 1
		;
	BEGIN TRANSACTION [T1]
		UPDATE [User]
			SET [DisplayName] = @DisplayName
			WHERE [Uuid] = @UserUuid
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
-- UpdateUserProfileBio --
-----------------------------------------------------------

-- USP_UpdateUserProfileBio replaces the existing bio with the provided one.
-- Parameters
--	@UserUuid:	UNIQUEIDENTIFIER the UserUuid for the user who's information should be updated.
--				Must be a valid v4 UUID.
--	@Bio NVARCHAR(1000) the bio to be added.
-- Outputs
--	Query result indicating the number of rows updated (should be 1)
-- Errors
--	50101: The provided UserUuid was null.
--	50301: No user found with the provided UserUuid.
CREATE PROCEDURE [USP_UpdateUserProfileBio]
	@UserUuid UNIQUEIDENTIFIER
	,@Bio NVARCHAR(1000)
AS
SET NOCOUNT ON
;
BEGIN TRY
	IF @UserUuid IS NULL
		THROW 50101, N'uuid must not be null', 1
		;
	IF NOT EXISTS (SELECT [Uuid] FROM [User] WHERE [Uuid] = @UserUuid)
		THROW 50301, N'user does not exist', 1
		;
	DECLARE @ProfileUuid UNIQUEIDENTIFIER
	SET @ProfileUuid = (
			SELECT [P].[Uuid] FROM [User] AS [U]
				INNER JOIN [UserProfile] AS [UP]
					ON [U].[Uuid] = [UP].[User_Uuid]
				INNER JOIN [Profile] AS [P]
					ON [UP].[Profile_Uuid] = [P].[Uuid]
				WHERE [U].[Uuid] = @UserUuid
		)
	BEGIN TRANSACTION [T1]
		UPDATE [Profile]
			SET [Bio] = @Bio
			WHERE [Uuid] = @ProfileUuid
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
-- UpdateUserProfileGravitarUrl --
-----------------------------------------------------------

-- USP_UpdateUserProfileGravitarUrl replaces the existing gravitar url with the provided one.
-- Parameters
--	@UserUuid:	UNIQUEIDENTIFIER the UserUuid for the user who's information should be updated.
--				Must be a valid v4 UUID.
--	@GravitarUrl NVARCHAR(1000) the gravitar url to be added.
-- Outputs
--	Query result indicating the number of rows updated (should be 1)
-- Errors
--	50101: The provided UserUuid was null.
--	50301: No user found with the provided UserUuid.
CREATE PROCEDURE [USP_UpdateUserProfileGravitarUrl]
	@UserUuid UNIQUEIDENTIFIER
	,@GravitarUrl NVARCHAR(1000)
AS
SET NOCOUNT ON
;
BEGIN TRY
	IF @UserUuid IS NULL
		THROW 50101, N'uuid must not be null', 1
		;
	IF NOT EXISTS (SELECT [Uuid] FROM [User] WHERE [Uuid] = @UserUuid)
		THROW 50301, N'user does not exist', 1
		;
	DECLARE @ProfileUuid UNIQUEIDENTIFIER
	SET @ProfileUuid = (
			SELECT [P].[Uuid] FROM [User] AS [U]
				INNER JOIN [UserProfile] AS [UP]
					ON [U].[Uuid] = [UP].[User_Uuid]
				INNER JOIN [Profile] AS [P]
					ON [UP].[Profile_Uuid] = [P].[Uuid]
				WHERE [U].[Uuid] = @UserUuid
		)
	BEGIN TRANSACTION [T1]
		UPDATE [Profile]
			SET [GravitarUrl] = @GravitarUrl
			WHERE [Uuid] = @ProfileUuid
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
-- UpdateUserProfileSharingGravitarUrl --
-----------------------------------------------------------

-- USP_UpdateUserProfileSharingGravitarUrl updates the sharing preference for the gravitar with other users.
-- Parameters
--	@UserUuid:	UNIQUEIDENTIFIER the UserUuid for the user who's information should be updated.
--				Must be a valid v4 UUID.
--	@Share: NCHAR(1) if gravitar should be shared with all users set 'Y', if not 'N'
-- Outputs
--	Query result indicating the number of rows updated (should be 1)
-- Errors
--	50101: The provided UserUuid was null.
--	50102: The provided Share was null.
--	50201: The provided Share was not one of 'Y' or 'N'
--	50301: No user found with the provided UserUuid.
CREATE PROCEDURE [USP_UpdateUserProfileSharingGravitarUrl]
	@UserUuid UNIQUEIDENTIFIER
	,@Share NCHAR(1)
AS
SET NOCOUNT ON
;
BEGIN TRY
	IF @UserUuid IS NULL
		THROW 50101, N'uuid must not be null', 1
		;
	IF @Share IS NULL
		THROW 50102, N'share must not be null', 1
		;
	IF @Share != N'Y' OR @Share != N'N'
		THROW 50102, N'share can only be Y or N', 1
		;
	IF NOT EXISTS (SELECT [Uuid] FROM [User] WHERE [Uuid] = @UserUuid)
		THROW 50301, N'user does not exist', 1
		;
	DECLARE @ProfileSharingUuid UNIQUEIDENTIFIER
	SET @ProfileSharingUuid = (
			SELECT [PS].[Uuid] FROM [User] AS [U]
				INNER JOIN [UserProfileSharing] AS [UPS]
					ON [U].[Uuid] = [UPS].[User_Uuid]
				INNER JOIN [ProfileSharing] AS [PS]
					ON [UPS].[ProfileSharing_Uuid] = [PS].[Uuid]
				WHERE [U].[Uuid] = @UserUuid
		)
	BEGIN TRANSACTION [T1]
		UPDATE [ProfileSharing]
			SET [GravitarUrl] = @Share
			WHERE [Uuid] = @ProfileSharingUuid
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
-- UpdateUserProfileSharingBio --
-----------------------------------------------------------

-- USP_UpdateUserProfileSharingBio updates the sharing preference for the bio with other users.
-- Parameters
--	@UserUuid:	UNIQUEIDENTIFIER the UserUuid for the user who's information should be updated.
--				Must be a valid v4 UUID.
--	@Share: NCHAR(1) if bio should be shared with all users set 'Y', if not 'N'
-- Outputs
--	Query result indicating the number of rows updated (should be 1)
-- Errors
--	50101: The provided UserUuid was null.
--	50102: The provided Share was null.
--	50201: The provided Share was not one of 'Y' or 'N'
--	50301: No user found with the provided UserUuid.
CREATE PROCEDURE [USP_UpdateUserProfileSharingBio]
	@UserUuid UNIQUEIDENTIFIER
	,@Share NCHAR(1)
AS
SET NOCOUNT ON
;
BEGIN TRY
	IF @UserUuid IS NULL
		THROW 50101, N'uuid must not be null', 1
		;
	IF @Share IS NULL
		THROW 50102, N'share must not be null', 1
		;
	IF @Share != N'Y' OR @Share != N'N'
		THROW 50102, N'share can only be Y or N', 1
		;
	IF NOT EXISTS (SELECT [Uuid] FROM [User] WHERE [Uuid] = @UserUuid)
		THROW 50301, N'user does not exist', 1
		;
	DECLARE @ProfileSharingUuid UNIQUEIDENTIFIER
	SET @ProfileSharingUuid = (
			SELECT [PS].[Uuid] FROM [User] AS [U]
				INNER JOIN [UserProfileSharing] AS [UPS]
					ON [U].[Uuid] = [UPS].[User_Uuid]
				INNER JOIN [ProfileSharing] AS [PS]
					ON [UPS].[ProfileSharing_Uuid] = [PS].[Uuid]
				WHERE [U].[Uuid] = @UserUuid
		)
	BEGIN TRANSACTION [T1]
		UPDATE [ProfileSharing]
			SET [Bio] = @Share
			WHERE [Uuid] = @ProfileSharingUuid
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
-- UpdateUserProfileSharingDisplayName --
-----------------------------------------------------------

-- USP_UpdateUserProfileSharingDisplayName updates the sharing preference for the display name with other users.
-- Parameters
--	@UserUuid:	UNIQUEIDENTIFIER the UserUuid for the user who's information should be updated.
--				Must be a valid v4 UUID.
--	@Share: NCHAR(1) if display name should be shared with all users set 'Y', if not 'N'
-- Outputs
--	Query result indicating the number of rows updated (should be 1)
-- Errors
--	50101: The provided UserUuid was null.
--	50102: The provided Share was null.
--	50201: The provided Share was not one of 'Y' or 'N'
--	50301: No user found with the provided UserUuid.
CREATE PROCEDURE [USP_UpdateUserProfileSharingDisplayName]
	@UserUuid UNIQUEIDENTIFIER
	,@Share NCHAR(1)
AS
SET NOCOUNT ON
;
BEGIN TRY
	IF @UserUuid IS NULL
		THROW 50101, N'uuid must not be null', 1
		;
	IF @Share IS NULL
		THROW 50102, N'share must not be null', 1
		;
	IF @Share != N'Y' OR @Share != N'N'
		THROW 50102, N'share can only be Y or N', 1
		;
	IF NOT EXISTS (SELECT [Uuid] FROM [User] WHERE [Uuid] = @UserUuid)
		THROW 50301, N'user does not exist', 1
		;
	DECLARE @ProfileSharingUuid UNIQUEIDENTIFIER
	SET @ProfileSharingUuid = (
			SELECT [PS].[Uuid] FROM [User] AS [U]
				INNER JOIN [UserProfileSharing] AS [UPS]
					ON [U].[Uuid] = [UPS].[User_Uuid]
				INNER JOIN [ProfileSharing] AS [PS]
					ON [UPS].[ProfileSharing_Uuid] = [PS].[Uuid]
				WHERE [U].[Uuid] = @UserUuid
		)
	BEGIN TRANSACTION [T1]
		UPDATE [ProfileSharing]
			SET [DisplayName] = @Share
			WHERE [Uuid] = @ProfileSharingUuid
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
-- UpdateUserSessionExpired --
-----------------------------------------------------------

-- USP_UpdateUserSessionExpired marks the provided session as expired.
-- Parameters
--	@SessionUuid UNIQUEIDENTIFIER the session to be expired.
-- Outputs
--	Query result indicating the number of rows updated (should be 1)
-- Errors
--	50101: The provided SessionUuid was null.
--	50301: No session found with the provided SessionUuid.
CREATE PROCEDURE [USP_UpdateUserSessionExpired]
	@SessionUuid UNIQUEIDENTIFIER
AS
SET NOCOUNT ON
;
BEGIN TRY
	IF @SessionUuid IS NULL
		THROW 50101, N'session uuid must not be null', 1
		;
	IF NOT EXISTS (SELECT [Uuid] FROM [Session] WHERE [Uuid] = @SessionUuid)
		THROW 50301, N'session does not exist', 1
		;
	BEGIN TRANSACTION [T1]
		UPDATE [Session]
			SET [Status] = N'Expired'
			WHERE [Uuid] = @SessionUuid
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
END CATCH
;
GO

----------------------------------------------------------------
-------- DELETE Procedures --------
----------------------------------------------------------------

-----------------------------------------------------------
-- DeleteUserEmail --
-----------------------------------------------------------

-- USP_DeleteUserEmail deletes the email from the list of emails for the given user.
-- Parameters
--	@UserUuid:	UNIQUEIDENTIFIER the UserUuid for the user who's information should be deleted.
--				Must be a valid v4 UUID.
--	@Email NVARCHAR(255) the email to be deleted.
-- Outputs
--	Query result indicating the number of rows updated (should be 1)
-- Errors
--	50101: The provided UserUuid was null.
--	50102: The provided Email was null.
--	50301: No user found with the provided UserUuid.
--	50302: Provided email not found in list of user's emails.
CREATE PROCEDURE [USP_DeleteUserEmail]
	@UserUuid UNIQUEIDENTIFIER
	,@Email NVARCHAR(255)
AS
SET NOCOUNT ON
;
BEGIN TRY
	IF @UserUuid IS NULL
		THROW 50101, N'uuid must not be null', 1
		;
	IF @Email IS NULL
		THROW 50102, N'email must not be null', 1
		;
	IF NOT EXISTS (SELECT [Uuid] FROM [User] WHERE [Uuid] = @UserUuid)
		THROW 50301, N'user does not exist', 1
		;
	IF NOT EXISTS (
			SELECT [U].[Uuid] 
				FROM [User] AS [U]
				INNER JOIN [UserEmail] AS [UE]
					ON [U].[Uuid] = [UE].[User_Uuid]
				INNER JOIN [Email] AS [E]
					ON [UE].[Email_Uuid] = [E].[Uuid]
				WHERE [U].[Uuid] = @UserUuid AND [E].[Email] = @Email
		)
		THROW 50302, N'email does not exist for user', 1
		;
	;
	DECLARE @EmailUuid UNIQUEIDENTIFIER
	;
	SET @EmailUuid = (
			SELECT [E].[Uuid] FROM [User] AS [U]
				INNER JOIN [UserEmail] AS [UE]
					ON [U].[Uuid] = [UE].[User_UUID]
				INNER JOIN [Email] AS [E]
					ON [UE].[Email_Uuid] = [E].[Uuid]
				WHERE [U].[Uuid] = @UserUuid AND [E].[Email] = @Email
		)
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
END CATCH
;
GO

-----------------------------------------------------------
-- DeleteUser --
-----------------------------------------------------------

-- USP_DeleteUser deletes the user and their associated data in the database.
-- Parameters
--	@UserUuid:	UNIQUEIDENTIFIER the UserUuid for the user who's information should be deleted.
--				Must be a valid v4 UUID.
-- Outputs
--	Query result indicating the number of rows updated (should be 1)
-- Errors
--	50101: The provided UserUuid was null.
--	50301: No user found with the provided UserUuid.
CREATE PROCEDURE [USP_DeleteUser]
	@UserUuid UNIQUEIDENTIFIER
AS
SET NOCOUNT ON
;
BEGIN TRY
	IF @UserUuid IS NULL
		THROW 50101, N'uuid must not be null', 1
		;
	IF NOT EXISTS (SELECT [Uuid] FROM [User] WHERE [Uuid] = @UserUuid)
		THROW 50301, N'user does not exist', 1
		;
	BEGIN TRANSACTION [T1]
		DELETE FROM [User]
			WHERE [Uuid] = @UserUuid
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
END CATCH
;
GO

-----------------------------------------------------------
-- DeleteUserSession --
-----------------------------------------------------------

-- USP_DeleteUserSession deletes the session from the list of sessions for the given user.
-- Parameters
--	@UserUuid:	UNIQUEIDENTIFIER the UserUuid for the user who's information should be deleted.
--				Must be a valid v4 UUID.
--	@SessionUuid UNIQUEIDENTIFIER the session to be deleted.
-- Outputs
--	Query result indicating the number of rows updated (should be 1)
-- Errors
--	50101: The provided UserUuid was null.
--	50102: The provided SessionUuid was null.
--	50301: No user found with the provided UserUuid.
--	50302: Provided session not found in list of user's sessions.
CREATE PROCEDURE [USP_DeleteUserSession]
	@UserUuid UNIQUEIDENTIFIER
	,@SessionUuid NVARCHAR(255)
AS
SET NOCOUNT ON
;
BEGIN TRY
	IF @UserUuid IS NULL
		THROW 50101, N'uuid must not be null', 1
		;
	IF @SessionUuid IS NULL
		THROW 50102, N'session uuid must not be null', 1
		;
	IF NOT EXISTS (SELECT [Uuid] FROM [User] WHERE [Uuid] = @UserUuid)
		THROW 50301, N'user does not exist', 1
		;
	IF NOT EXISTS (SELECT [Uuid] FROM [Session] WHERE [Uuid] = @SessionUuid)
		THROW 50302, N'session does not exist', 1
		;

	BEGIN TRANSACTION [T1]
		DELETE FROM [Session]
			WHERE [Uuid] = @SessionUuid
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
END CATCH
;
GO