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

USE [Perceptia]
;
GO

-------------------------------------------------------------------------------
-- Examples --
-------------------------------------------------------------------------------

-----------------------------------------------------------
-- TODO --
-----------------------------------------------------------
/*
	Schema Version: 0.3.2
	Procedure Version: 0.5.0
*/
-- TODO: Make this into an example

-- Run this first
DECLARE @UserUuid UNIQUEIDENTIFIER
SET @UserUuid = NEWID()
;
EXECUTE USP_CreateUser 
	@Uuid = @UserUuid, 
	@UserName = N'test',
	@FullName = N'test fn',
	@DisplayName = N'TEST',
	@EncodedHash = N'NOT A REAL HASH'
; 
GO

-- Copy uuid created in last usp to this
EXECUTE USP_AddUserEmail
	@Uuid = N'557274CD-30BA-4EA2-81F4-286BFB269F5B',
	@Email = 'testing@test.com'
;
GO

-- Copy uuid created in first usp to this
EXECUTE USP_GetUserEmailByUuid
	@Uuid = N'557274CD-30BA-4EA2-81F4-286BFB269F5B'
;
GO

-- Shows emails in Email table
SELECT TOP (1000) [Uuid]
      ,[Email]
      ,[Created]
  FROM [Perceptia].[dbo].[Email]