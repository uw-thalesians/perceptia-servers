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

    public static function listAllQuizzesJSON($root) {
        header('Content-type: application/json');

        $ch = curl_init();

        $umap_py = "localhost/py/umap.py?root=" . urlencode($root);

        curl_setopt($ch, CURLOPT_URL, $umap_py);
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);

        $resp = curl_exec($ch);

        curl_close($ch);

        $topics = array();

        foreach( $resp["topic"] as $topic ) {
            $topics[] = array(
                                "keyword"=>$topic->keyword,
                                "url"=>$topic->url,
                                "x"=>$topic->x,
                                "y"=>$topic->y);
        }

        echo json_encode(array("rest_api_v" => "1.1", "root"=>$root, "topics"=> $topics));
    }
}

?>
