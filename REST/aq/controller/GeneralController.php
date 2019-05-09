<?php

class GeneralController
{    
    public static function listAllQuizzesJSON()
    {
        $f3 = Base::instance();

        $conn = new Connection();

        $quizzes = $conn->getAllQuizzes();

        header('Content-type: application/json');

        $quiz_list = array();
        foreach( $quizzes as $quiz ) {
            $quiz_list[] = array("keyword" =>$quiz->keyword);
        }
        echo json_encode(array("quizzes" => $quiz_list));
    }
}

?>
