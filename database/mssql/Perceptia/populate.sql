/*
	Title: Perceptia Database Populate
	Version: 0.1.0
*/
-------------------------------------------------------------------------------
-- Change Log --
-------------------------------------------------------------------------------
/*
	Date, Changer, Short Description, Version
	2019/05/19, Chris, Created Populate, 0.1.0
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

USE [Perceptia]
;
GO

-------------------------------------------------------------------------------
-- Populate Database --
-------------------------------------------------------------------------------

-----------------------------------------------------------
-- Version Table --
-----------------------------------------------------------

INSERT INTO [Version] ([Uuid], [Name], [Version], [Description])
VALUES (N'CE8A00FF-5715-424C-B313-28E280F7165B', N'Populate', N'0.1.0', N'The Perceptia Database Populate.')
;
GO