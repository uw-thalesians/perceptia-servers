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

        echo json_encode(
                        array(  
                                "rest_api_v"=> "1.1",
                                "questions"=>$quiz->questions,
                                "paragraphs" => $quiz->paras
                            )
                        );
    }

    public static function editQuestionJSON($user)
    {
        $conn = new Connection();
        
        $json = file_get_contents('php://input');
        
        $obj = json_decode($json, true);
        $questionID = $obj['questionID'];
        
        //print_r($questionID);

        $newText = $obj['newText'];

        //print_r($newText);

        $newAnswer = $obj['newAnswer'];

        //print_r($newAnswer);

        $result = $conn->editQuestion($user, $questionID, $newText, $newAnswer);

        echo json_encode(array("questionID" => $questionID, "result" => $result));
    }
    
    public static function deleteQuestionJSON($user)
    {
        $conn = new Connection();
        
        $json = file_get_contents('php://input');
        
        $obj = json_decode($json, true);
        $questionID = $obj['questionID'];

        $result = $conn->deleteQuestion($user, $questionID);

        echo json_encode(array("questionID" => $questionID, "result" => $result));
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
