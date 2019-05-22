DROP DATABASE IF EXISTS `any_quiz_db`;

CREATE DATABASE `any_quiz_db`;

USE `any_quiz_db`;

CREATE TABLE `any_quiz_db`.`quiz_questions` (
    `question` BLOB NOT NULL ,
    `id` INT NOT NULL AUTO_INCREMENT ,
    `answer` VARCHAR(30) NOT NULL ,
    `quiz_id` INT NOT NULL,
    `q_type` INT NOT NULL DEFAULT '1',
    `p_id` INT NOT NULL,
    /*FOREIGN KEY(quiz_id)
      REFERENCES quizzes(id),*/
    PRIMARY KEY(id)
);

CREATE TABLE `any_quiz_db`.`quizzes` (
    `keyword` VARCHAR(255) NOT NULL,
    `id` INT NOT NULL AUTO_INCREMENT,
    `image` VARCHAR(2084),
    `when` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `source` VARCHAR(64) NOT NULL DEFAULT 'wiki',
    `total_read_count` INT NOT NULL DEFAULT '1', 

    `status` INT NOT NULL DEFAULT '0',
    PRIMARY KEY(id)
);

CREATE TABLE `any_quiz_db`.`paragraphs` (
    `quiz_id` VARCHAR(255) NOT NULL,
    `id` INT NOT NULL AUTO_INCREMENT,
    `text` BLOB NOT NULL,
    PRIMARY KEY(id)
);