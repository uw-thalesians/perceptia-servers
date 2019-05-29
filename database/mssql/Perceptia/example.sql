/*
	Title: Perceptia Examples
*/

-------------------------------------------------------------------------------
-- Sections --
-------------------------------------------------------------------------------

-- Setup
-- Examples

-------------------------------------------------------------------------------
-- Setup --
-------------------------------------------------------------------------------

/*
USE [Perceptia]
;
GO
*/

-------------------------------------------------------------------------------
-- Examples --
-------------------------------------------------------------------------------

-----------------------------------------------------------
-- TODO --
-----------------------------------------------------------
/*
	Schema Version: 0.7.1
	Procedure Version: 0.8.1
*/
-- TODO: Make this into an example

-- Run this first
DECLARE @UserUuid UNIQUEIDENTIFIER
SET @UserUuid = N'557274CD-30BA-4EA2-81F4-286BFB269F5B'
;
EXECUTE USP_CreateUser 
	@UserUuid = @UserUuid, 
	@UserName = N'test',
	@FullName = N'test fn',
	@DisplayName = N'TEST',
	@EncodedHash = N'NOT A REAL HASH'
; 
GO

-- Copy uuid created in last usp to this
EXECUTE USP_CreateUserEmail
	@UserUuid = N'557274CD-30BA-4EA2-81F4-286BFB269F5B',
	@Email = 'testing@test.com'
;
GO

-- Copy uuid created in first usp to this
EXECUTE USP_ReadUserEmailsByUuid
	@UserUuid = N'557274CD-30BA-4EA2-81F4-286BFB269F5B'
;
GO

-- Shows emails in Email table
SELECT TOP (1000) [Uuid]
      ,[Email]
      ,[Created]
  FROM [Perceptia].[dbo].[Email]

 -- Shows users
SELECT TOP (1000) [Uuid]
      ,[Username]
      ,[Created]
  FROM [Perceptia].[dbo].[User]



-- Copy uuid created in first usp to this
EXECUTE USP_UpdateUserProfileGravitarUrl
	@UserUuid = N'557274CD-30BA-4EA2-81F4-286BFB269F5B'
	,@GravitarUrl = N'https://s.gravatar.com/avatar/9d6f7ef1d8a8e975c66d2b47e0e45402?s=80'
;
GO

-- Copy uuid created in first usp to this
EXECUTE USP_UpdateUserProfileBio
	@UserUuid = N'557274CD-30BA-4EA2-81F4-286BFB269F5B'
	,@Bio = N'My name is name'
;
GO

-- Copy uuid created in first usp to this
EXECUTE USP_UpdateUserProfileSharingBio
	@UserUuid = N'557274CD-30BA-4EA2-81F4-286BFB269F5B'
	,@Share = N'Y'
;
GO



-- Copy uuid created in first usp to this
EXECUTE USP_ReadUserProfileByUuid
	@UserUuid = N'557274CD-30BA-4EA2-81F4-286BFB269F5B'
;
GO

-- Copy uuid created in first usp to this
EXECUTE USP_UpdateUserUsername
	@UserUuid = N'557274CD-30BA-4EA2-81F4-286BFB269F5B'
	,@Username = N'test'
;
GO

-- Copy uuid created in first usp to this
EXECUTE USP_DeleteUser
	@UserUuid = N'557274CD-30BA-4EA2-81F4-286BFB269F5B'
;
GO