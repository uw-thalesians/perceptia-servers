/*
	Title: Perceptia Database Populate
	Version: 0.2.0
*/
-------------------------------------------------------------------------------
-- Change Log --
-------------------------------------------------------------------------------
/*
	Date, Changer, Short Description, Version
	2019/05/19, Chris, Created Populate, 0.1.0
	2019/05/21, Chris, Update versions for schema and proc, 0.2.0
*/

-------------------------------------------------------------------------------
-- Sections --
-------------------------------------------------------------------------------

-- Setup
-- Populate Database
-- Create Procedures

-------------------------------------------------------------------------------
-- Setup --
-------------------------------------------------------------------------------

/*
USE [Perceptia]
;
GO
*/

-------------------------------------------------------------------------------
-- Populate Database --
-------------------------------------------------------------------------------

-----------------------------------------------------------
-- Procedure Version Table --
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
		,N'1.0.0'
		,N'The Perceptia Database Stored Procedures.'
	)
;
GO

-----------------------------------------------------------
-- Schema Version Table --
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
		,N'1.0.0'
		,N'The Perceptia Database Schema.'
	)
;
GO

-----------------------------------------------------------
-- Version Table --
-----------------------------------------------------------

INSERT INTO [Version] ([Uuid], [Name], [Version], [Description])
VALUES (N'CE8A00FF-5715-424C-B313-28E280F7165B', N'Populate', N'0.2.0', N'The Perceptia Database Populate.')
;
GO