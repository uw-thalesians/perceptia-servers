<?php

ini_set('max_execution_time', 120);

session_start();

$f3 = require('vendor/bcosca/fatfree-core/base.php');

$f3->set('ONERROR', function($f3){
	echo $f3->get('ERROR.code');
	echo $f3->get('ERROR.status');
	echo $f3->get('ERROR.text');
	echo $f3->get('ERROR.trace');
	echo $f3->get('ERROR.level');

} );

$f3->set("CORS.origin", "*");

$f3->set('DEBUG', 3);

$f3->set('AUTOLOAD', 'sql/;view/;controller/;model/;py/;');

$f3->route('GET /v1/read/@keyword', function($f3){
    $quizController = new QuizController();
    $keyword = $f3->get('PARAMS.keyword');

    header('Content-Type: application/json');
    $quizController->startJSON($keyword);
});

$f3->route('GET /v1/list', 'GeneralController::listAllQuizzesJSON');

$f3->route('GET /v1/questions/@keyword', function($f3) {

    $quizController = new QuizController();

    $keyword = $f3->get('PARAMS.keyword');

    header('Content-Type: application/json');
    $quizController->questionsJSON($keyword);

});

$f3->route('POST /v1/questions/grade', function() {
    $quizController = new QuizController();

    header('Content-Type: application/json');
    $quizController->gradeJSON($_SESSION['user']);
});

$f3->run();
