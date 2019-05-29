<?php

class GeneralController
{    
    public static function listAllQuizzesJSON($start, $end, $sort) {

        if(!isset($start)) {
            $start = 0;
        }
        
        if(!isset($end)) {
            $end = -1;
        }

        $f3 = Base::instance();

        $conn = new Connection();

        $results = $conn->getAllQuizzes($start, $end, $sort);

        $quizzes = $results["quizzes"];

        header('Content-type: application/json');

        $quiz_list = array();
        foreach( $quizzes as $quiz ) {
            $quiz_list[] = array("timestamp"=>$quiz->timestamp, "keyword" =>$quiz->keyword, "source"=>$quiz->source, "image"=>$quiz->image);
        }

        echo json_encode(array("rest_api_v"=> "1.1", "quizzes" => $quiz_list, "sort"=>$results["sort"], "start"=>$results["start"], "end"=>$results["end"]));
    }

    public static function quizStatusJSON($keyword, $source) {

        $conn = new Connection();

        $results = $conn->getQuizStatus($keyword, $source);

        $resp = array("rest_api_v"=> "1.1",
                         "keyword" => $keyword,
                         "source" => $source);

        foreach($results as $key=>$val) {
            $resp[$key]=$val;
        }

        echo json_encode($resp);
    }
}

?>
