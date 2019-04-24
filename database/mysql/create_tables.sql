CREATE DATABASE `any_quiz_db`;

USE `any_quiz_db`;

CREATE TABLE `any_quiz_db`.`quiz_users` (
    `username` VARCHAR(20) NOT NULL ,
    `id` INT NOT NULL AUTO_INCREMENT ,
    `totalscore` INT NOT NULL ,
    `streak` INT NOT NULL ,
	`best_streak` INT NOT NULL DEFAULT '0',
	`profilePicture` VARCHAR(255) NOT NULL,
	`password` VARCHAR(255) NOT NULL,
	PRIMARY KEY(id)
);

CREATE TABLE `any_quiz_db`.`quiz_questions` (
    `question` BLOB NOT NULL ,
    `id` INT NOT NULL AUTO_INCREMENT ,
    `answer` VARCHAR(30) NOT NULL ,
    `quiz_id` INT NOT NULL,
    /*FOREIGN KEY(quiz_id)
      REFERENCES quizzes(id),*/

	PRIMARY KEY(id)
);

CREATE TABLE `any_quiz_db`.`quizzes` (
    `keyword` VARCHAR(255) NOT NULL,
    `summary` BLOB NOT NULL,
    `id` INT NOT NULL AUTO_INCREMENT,
	`image` VARCHAR(256) NOT NULL,
 	PRIMARY KEY(id)
);