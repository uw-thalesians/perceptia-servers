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

$f3->route('GET /api/v1/anyquiz/read/@keyword', function($f3) {
    
    
    header('Content-Type: application/json');

    $source = "wiki";

    if( $f3->exists('GET.source') ) {
        $source = $f3->get('GET.source');
    }

    $quizController = new QuizController();
    $keyword = $f3->get('PARAMS.keyword');
    $quizController->startJSON( $keyword, $source );
});

$f3->route('GET /api/v1/anyquiz/list', function($f3) {
    $start = 0;
    $end = -1;
    $sort = 'alpha';

    if($f3->exists('GET.start')) {
        $start = $f3->get('GET.start');
    }
    
    if($f3->exists('GET.end')) {
        $end = $f3->get('GET.end');
    }

    if($f3->exists('GET.sort')) {
        $sort = $f3->get('GET.sort');
    }

    GeneralController::listAllQuizzesJSON($start, $end, $sort);

});

$f3->route('GET /api/v1/anyquiz/questions/@keyword', function($f3) {

    $quizController = new QuizController();

    $source = 'wiki';

    if( $f3->exists('GET.source') ) {
        $source = $f3->get('GET.source');
    }

    $keyword = $f3->get('PARAMS.keyword');

    header('Content-Type: application/json');
    $quizController->questionsJSON($keyword, $source);

});

$f3->route('POST /api/v1/anyquiz/questions/grade', function() {
    header('Content-Type: application/json');
    
    $quizController = new QuizController();
    
    $quizController->gradeJSON($_SESSION['user']);
});

$f3->run();
