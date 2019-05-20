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

    public static function umapJSON($root) {
        header('Content-type: application/json');

        $ch = curl_init();

        $umap_py = "http://localhost/py/umap_conceptnet.py";//?root=" . urlencode($root);

        curl_setopt($ch, CURLOPT_URL, $umap_py);
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);

        $resp = curl_exec($ch);

        //echo $resp;
        $resp = json_decode($resp, true);

        curl_close($ch);

        $topics = array();

        foreach( $resp["umap"] as $topic=>$info ) {

            $topics[] = array(
                                "keyword"=>$topic,
                                "color"=>$info["color"],
                                "x"=>$info["x"],
                                "y"=>$info["y"]);
        }

        echo json_encode(array("rest_api_v" => "1.1", "root"=>$root, "topics"=> $topics));
    }
}

?>
