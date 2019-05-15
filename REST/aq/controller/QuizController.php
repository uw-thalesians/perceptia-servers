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
    
    public static function startJSON($keyword, $source)
    {
        $quiz = QuizController::$conn->findQuiz($keyword, $source);

        echo json_encode(array("rest_api_v" => "1.1","source" => $quiz->source, "summary" => $quiz->summary, "image" => $quiz->image, "timestamp" => $quiz->timestamp));
    }

    public static function questionsJSON($keyword, $source)
    {
        $quiz = QuizController::$conn->getQuizQuestions($keyword, $source);

        echo json_encode(array("questions"=>$quiz->questions));
    }

    public static function studyJSON($keyword, $source)
    {
        $quiz = QuizController::$conn->getStudyQuestions($keyword, $source);

        echo json_encode(array("questions"=>$quiz->questions, "paragraphs" => $quiz->paras));
    }

    public static function gradeJSON($user)
    {
        
        $conn = new Connection();
        
        $json = file_get_contents('php://input');
        
        $obj = json_decode($json, true);
        $questionID = $obj['questionID'];
        
        $boolResult = $conn->grade($user,
            $questionID,
            $obj['selectedAnswer']);

        echo json_encode(array("questionID" => $questionID, "result" => $boolResult));
    }
};
