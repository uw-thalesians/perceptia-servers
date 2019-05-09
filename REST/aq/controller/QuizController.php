<?php

class QuizController{

    private static $f3;
    private static $conn;

    public function __construct()
    {
        if(!isset(QuizController::$f3)) {
            QuizController::$f3 = Base::instance();
            QuizController::$conn = new Connection();
        }
    }
    
    public static function startJSON($keyword)
    {
        $quiz = QuizController::$conn->findQuiz($keyword, false);

        echo json_encode(array("summary"=>$quiz->summary));
    }

    public static function questionsJSON($keyword)
    {
        $quiz = QuizController::$conn->getQuizQuestions($keyword);

        echo json_encode(array("questions"=>$quiz->questions));
    }

    public static function gradeJSON($user)
    {
        $conn = new Connection();

        $json = file_get_contents('php://input');

        $obj = json_decode($json, true);
        $questionID = $obj['questionID'];

        $boolResult = $conn->grade($user,
            //$obj['questionID'],
            $questionID,
            $obj['selectedAnswer']);

        header('Content-type: application/json');

        echo json_encode(array("questionID" => $questionID, "result" => $boolResult));
    }
};