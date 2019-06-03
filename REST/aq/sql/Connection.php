<?php

require_once 'config.php';
require_once 'status.php';

class Connection
{
    

    private static $conn;

    private $_f3;

	public function __construct()
	{
	    $this->_f3 = Base::instance();

		//signleton, is this threadsafe in php? (it is in MPM prefork, but not with the multithreaded apache/php build)
		if(!isset(Connection::$conn))
		{
			try{
				Connection::$conn = new PDO(DB_DSN, DB_USERNAME, DB_PASSWORD, array(PDO::MYSQL_ATTR_INIT_COMMAND => "SET NAMES 'utf8'"));
				Connection::$conn->setAttribute( PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION );
			}catch (Exception $e) {
				echo "Connection failed " . $e->getMessage();
			}
        }
        
	}
	
	private function getConnection()
	{
		return new DB\SQL(
			DB_DSN,
			DB_USERNAME,
			DB_PASSWORD
		);
	}

    public function findQuiz($keyword, $source)
    {

        if(!isset($source)) {
            $source = "wiki";
        }

        $sql = "SELECT * FROM quizzes WHERE keyword=:keyword and source=:source";

        $quiz = null;

        try {

            // print_r($keyword);
            // print_r($source);

            $stmt = Connection::$conn->prepare($sql);

            $stmt->bindValue(":keyword", $keyword, PDO::PARAM_STR);
            $stmt->bindValue(":source", $source, PDO::PARAM_STR);

            $stmt->execute();

            $row = $stmt->fetch(PDO::FETCH_ASSOC);

            if($row) {

                $update_read = "UPDATE quizzes SET total_read_count = total_read_count + 1 WHERE id=:quiz_id";

                $select_para = "SELECT id, text FROM paragraphs WHERE quiz_id=:quiz_id";

                $id = $row["id"];

                $stmt = Connection::$conn->prepare($update_read);
                $stmt->bindValue(":quiz_id", $id, PDO::PARAM_INT);
                $stmt->execute();

                $stmt = Connection::$conn->prepare($select_para);
                $stmt->bindValue(":quiz_id", $id, PDO::PARAM_INT);
                $stmt->execute();

                $paras = $stmt->fetchAll(PDO::FETCH_ASSOC);

                $row['paras'] = $paras;

                //print_r($row);
                $quiz = new Quiz( $row );

            } else {
                $quiz_id = $this->storeQuizRequest($keyword, $source);
                // print_r("row not found");
                $summary = $this->fetchSummary($keyword, $source);
                // print_r($summary);
                $quiz = $this->addNewQuiz($keyword, $summary, $source, $quiz_id);
            }

        } catch (Exception $e) {
            echo "Error finding quiz: " . $e->getMessage();
        }

        return $quiz;
    }

    private function storeQuizRequest($keyword, $source) {
        $sql = "INSERT INTO quizzes (keyword, image, source) VALUES (:keyword, :image, :source)";

        $quiz_id = null;
        try {

            $stmt = Connection::$conn->prepare($sql);

            $stmt->bindValue(":keyword", $keyword);
            $stmt->bindValue(":image", "");
            $stmt->bindValue(":source", $source);
            //default value of 0 status inserted

            $stmt->execute();
          
            $quiz_id = Connection::$conn->lastInsertId();

        } catch (Exception $e){
            echo "Error creating quiz: " . $e->getMessage();
        }
      
      return $quiz_id;
    }

    private function fetchSummary($keyword, $source) {

        $curl_handle= curl_init();
        $summary = "";

        $safe_keyword = urlencode($keyword);

        switch($source) {
            case "solr_url":
                
                
                $solr_stream_url = "http://aqsolr:8983/solr/stream/update/extract?uprefix=attr_&fmap.content=body&commit=true";
                // print_r($solr_stream_url);
                
                curl_setopt($curl_handle, CURLOPT_URL, $solr_stream_url);
                curl_setopt($curl_handle, CURLOPT_POST, 1);
                curl_setopt($curl_handle, CURLOPT_POSTFIELDS, "stream.url=" . $safe_keyword);
                curl_setopt($curl_handle, CURLOPT_RETURNTRANSFER, 1);
                //curl_setopt($curl_handle, CURLOPT_PORT, 8983);

                $json = curl_exec($curl_handle);

                // print_r($json);
                // print_r(curl_error($curl_handle));
                // print_r(curl_getinfo($curl_handle, CURLINFO_HTTP_CODE));
                
                $solr_select_stream_name = "http://aqsolr:8983/solr/stream/select?q=attr_stream_name:" . "\"". $safe_keyword . "\"";
                // print_r($solr_select_stream_name);

                curl_setopt($curl_handle, CURLOPT_URL, $solr_select_stream_name);
                curl_setopt($curl_handle, CURLOPT_POST, 0);

                $response = curl_exec($curl_handle);
                
                // print_r($response);
                // print_r(curl_error($curl_handle));
                // print_r(curl_getinfo($curl_handle, CURLINFO_HTTP_CODE));

                $response = json_decode($response, true)["response"];
                // print_r($response);

                //print_r($response["docs"][0]["attr_body"][0]);
//https://github.com/commonsense/conceptnet5/wiki/API
                switch($response["docs"][0]["attr_stream_content_type"])
                {
                    case "text/html; charset=utf-8":
                        //print_r("using text/html; charset=utf-8 splitting strategy");
                        //check if it's possible to keep original html and use something like beautiful soup
                        //to extract p tags
                        $summary = explode("\n \n postPage", $response["docs"][0]["attr_body"][0]);
                        break;
                    case "application/pdf":
                        //print_r("using application/pdf splitting strategy");
                        $summary = explode("\n \n page", $response["docs"][0]["attr_body"][0]);
                        break;

                    default:
                        //print_r("using default splitting strategy");
                        $summary = explode("\n \n ", $response["docs"][0]["attr_body"][0]);
                }
                
                #print_r($summary);

                break;

            case "wiki":
                $server = "localhost";
                //curl call py script
                $wikipedia_python_rest_query = $server.$this->_f3->get("BASE")."/py/w_rest.py?keyword=" . urlencode($keyword);
                #$wikipedia_python_rest_query = $server . "/?keyword=" . urlencode($keyword);
                #print_r($wikipedia_python_rest_query);

                curl_setopt($curl_handle, CURLOPT_URL, $wikipedia_python_rest_query);
                curl_setopt($curl_handle, CURLOPT_RETURNTRANSFER, 1);
                #curl_setopt($curl_handle, CURLOPT_PORT, 27277);

                $json = curl_exec($curl_handle);
                #print_r(curl_error($curl_handle));
                #print_r(curl_getinfo($curl_handle, CURLINFO_HTTP_CODE));
                
                #print_r($json);

                $summary = json_decode($json, true)["summary_text"];

                break;
        }

        curl_close($curl_handle);


        $this->updateQuizStatus($quiz->keyword, $quiz->source, STATUS::RETRVD_MEDIA);
        #print_r($summary);

        return $summary;
    }

    private function updateQuizStatus($keyword, $source, $status) {
        $sql = "UPDATE quizzes SET status=:status WHERE keyword=:keyword and source=:source";

        try {

            $stmt = Connection::$conn->prepare($sql);
            
            $stmt->bindValue(":status", $status, PDO::PARAM_INT);
            $stmt->bindValue(":keyword", $keyword, PDO::PARAM_STR);
            $stmt->bindValue(":source", $source, PDO::PARAM_STR);

            $stmt->execute();

        } catch (Exception $e){
            echo "Error creating quiz: " . $e->getMessage();
        }
    }

    public function getQuizStatus($keyword, $source) {
        $sql = "SELECT status FROM quizzes WHERE keyword=:keyword and source=:source";

        $result = array();

        try {

            $stmt = Connection::$conn->prepare($sql);

            $stmt->bindValue(":keyword", $keyword, PDO::PARAM_STR);
            $stmt->bindValue(":source", $source, PDO::PARAM_STR);

            $stmt->execute();

            $rows = $stmt->fetchAll(PDO::FETCH_ASSOC);
            
            if(count($rows)==0) {
                $result["error"] = STATUS::STATUS_NOTFOUND . " with keyword=$keyword and source=$source";
            }else {
                $status = $rows[0]["status"];
                $result["progress"] = $status;
                $result["total_steps"] = STATUS::getCount();
                $result["status_string"] = STATUS::STATUS_STRINGS[$status];
            }
        } catch (Exception $e){
            echo "Error finding quiz: " . $e->getMessage();
            $result["error"] = STATUS::STATUS_NOTFOUND . " with keyword=$keyword and source=$source";
        }

        return $result;
    }

    public function deleteQuestion($user, $questionID) {

        $sql = "DELETE FROM quiz_questions WHERE id=:questionID";

        $result = array();

        try {

            $stmt = Connection::$conn->prepare($sql);
            
            //$stmt->bindValue(":owner", $user, PDO::PARAM_STR);
            $stmt->bindValue(":questionID", $questionID, PDO::PARAM_INT);

            $stmt->execute();

            //$rows = $stmt->fetchAll(PDO::FETCH_ASSOC);

            //print_r($rows);
            $result["status"] = "Delete Successful";

        } catch (Exception $e){
            
            $msg = $e->getMessage();
            //echo "Error finding question: " . $e->getMessage();
            $result["error"] = STATUS::STATUS_NOTFOUND . " with questionID=$questionID : $msg";

        }

        return $result;
    }

    public function editQuestion($user, $questionID, $newText, $newAnswer) {

        $select_sql = "SELECT answer from quiz_questions where id=:questionID";

        $sql = "UPDATE quiz_questions SET question=:newText, answer=:newAnswer WHERE id=:questionID";

        $result = array();

        try {

            $stmt = Connection::$conn->prepare($select_sql);

            $stmt->bindValue(":questionID", $questionID, PDO::PARAM_INT);
            
            $stmt->execute();

            $rows = $stmt->fetchAll(PDO::FETCH_ASSOC);

            $result["prevAnswer"] = $rows[0]["answer"];
            $result["newAnswer"] = $newAnswer;

            $stmt = Connection::$conn->prepare($sql);
            
            //$stmt->bindValue(":owner", $user, PDO::PARAM_STR);
            $stmt->bindValue(":questionID", $questionID, PDO::PARAM_INT);
            $stmt->bindValue(":newText", $newText, PDO::PARAM_STR);
            $stmt->bindValue(":newAnswer", $newAnswer, PDO::PARAM_STR);

            $stmt->execute();

            //$rows = $stmt->fetchAll(PDO::FETCH_ASSOC);

            //print_r($rows);

        } catch (Exception $e){
            
            $msg = $e->getMessage();
            //echo "Error finding question: " . $msg;
            $result["error"] = STATUS::STATUS_NOTFOUND . " with questionID=$questionID : $msg";
            
        }

        return $result;
    }

    private function addNewQuiz($keyword, $summary, $source, $quiz_id)//, $lang)
    {
        //https://cse.google.com/cse/create/new
        //https://developers.google.com/custom-search/json-api/v1/introduction#identify_your_application_to_google_with_api_key
        //https://console.developers.google.com/apis/credentials?project=my-project-1489383889696
        //https://developers.google.com/apis-explorer/#p/customsearch/v1/

        //print_r($curl);

        #$cmdline = "/opt/python27/bin/python2.7 -c 'from py import n; print n.get_noun(\"" . explode(".", $summary)[0] . "\")' 2>&1";

        //$nltk_rest_query = $this->_f3->get("BASE")+"/py/n_rest.py?get_noun"+

        //$curl_handle = curl_init();
        //curl_setopt($curl_handle, CURLOPT_URL, $nltk_rest_query);
        //curl_setopt($curl_handle, CURLOPT_RETURNTRANSFER, 1);
        //$json = curl_exec($curl_handle);
        //curl_close($curl_handle);

        //$handle = popen($cmdline, "r");

        //$additional_keyword = json_decode($json, 1)["noun"];// fread($handle, 4096);

        //print_r($additional_keyword);

        //fclose($handle);

        //$matches = array();
        //preg_match('/[a-z][A-Z]/', $additional_keyword, $matches);

        //$additional_keyword = $matches[1][0];
        //$additional_keyword = preg_replace('/[^A-Za-z0-9\-]/', '', $additional_keyword);

        //print_r($additional_keyword);

        $tries = 0;

        //do {

        $curl = curl_init();
        $search_keyword = str_replace(" ", "+", $keyword);// . ($tries==0 ? '+' . $additional_keyword : '');

        //print_r($search_keyword);

        //&rights=cc_publicdomain+cc_sharealike
        //&imgSize=large
        //&cx=004799634748936919555:ewzgppgp6wu
        $google_cse_rest_api_get = "https://www.googleapis.com/customsearch/v1?q=%22$search_keyword%22&searchType=image&safe=high&key=$GOOGLE_API_KEY";

        //print_r($google_cse_rest_api_get);

        //Google Custom Search v1 REST API
        curl_setopt($curl, CURLOPT_URL, $google_cse_rest_api_get);
        curl_setopt($curl, CURLOPT_RETURNTRANSFER, 1);

        $json = curl_exec($curl);

        //print_r($json);

        $search = json_decode($json, true);
        // $search = array();
        $search_results = 0;

        //var_dump($search);

        if(array_key_exists('queries', $search)) {
            $search_results = $search['queries']['request']['totalResults'];
        } else if(array_key_exists('error', $search)) {
            //print_r($search['error']);
            $search_results = 0;
        }

        $tries++;
        //echo $tries;
        ////curl_close($curl);
        //print_r($search_results);
        //print_r($tries);
        //} while($results == "0" || $tries < 2);

        $accepted_filetypes = array( "jpg", "png" );

        if(isset($search))
        {
            $curl = curl_init();
            //var_dump($search);

            #response header
            curl_setopt( $curl, CURLOPT_HEADER, true);
            curl_setopt( $curl, CURLOPT_RETURNTRANSFER, true);
            curl_setopt( $curl, CURLOPT_NOBODY, true);

            $tries = 0;

            do {

                $imageURL = $search['items'][$tries]['link'];
                $imageURL = explode("?", $imageURL)[0];
                //print_r($imageURL);

                $pathComponents = pathinfo($imageURL);

                //print_r($pathComponents);

                if(in_array($pathComponents['extension'], $accepted_filetypes)) {

                    $filename = $pathComponents['basename'];
                    $filename = str_replace(" ", "_", urldecode($filename));
                    //print_r($filename);

                    $remotePath = $imageURL;
                }
                $tries++;

                curl_setopt( $curl, CURLOPT_URL, $remotePath);

                $result = curl_exec( $curl );

                //$header = curl_getinfo($curl);
                
                $filesize = curl_getinfo($curl, CURLINFO_CONTENT_LENGTH_DOWNLOAD);
                
                //print_r($filesize);

            } while(($filesize < 10000) && $tries < 10);

            curl_close($curl);
        }else{
            //var_dump($json);
        }

        if(!isset($remotePath))
            $remotePath = 'images/quiz.png';

        $sql = "UPDATE quizzes SET image=:image where keyword=:keyword AND source=:source";

        try {

            $stmt = Connection::$conn->prepare($sql);

            $stmt->bindValue(":keyword", $keyword);
            $stmt->bindValue(":image", $remotePath);
            $stmt->bindValue(":source", $source);
            //default value of 0 status inserted

            $stmt->execute();

            //since bind value expects unique tokens, execute these one at a time for now, but look for
            //"the right way" i.e. [...] VALUES (),()[...] as time permits

            $sql = "insert into paragraphs (quiz_id, text) VALUES (:quiz_id, :text)";

            foreach ( $summary as $para ) {

                $stmt = Connection::$conn->prepare($sql);

                $stmt->bindValue(":quiz_id", $quiz_id, PDO::PARAM_INT);
                $stmt->bindValue(":text", $para, PDO::PARAM_STR);

                $stmt->execute();
            }

            $quiz = $this->findQuiz($keyword, $source);

        } catch (Exception $e){
            echo "Error creating quiz: " . $e->getMessage();
        }

        //in this older impl, the media was already retrieved before the first insert,
        //this will be here in newer code, so it is here now rather than changing the previous
        //insert and dealing with it again later in a merge conflict resolution
        $this->updateQuizStatus($quiz->keyword, $quiz->source, STATUS::QG_BEGINS);

        $path_to_py = dirname($_SERVER['PHP_SELF']) . "/py/n_cgi.py";


        #print_r($path_to_py);

        $nltk_py_path = "localhost${path_to_py}?keyword=" . urlencode($quiz->keyword);// . "&lang=" . urlencode($lang);

        #print_r($nltk_py_path);
        $curl = curl_init();
        curl_setopt($curl, CURLOPT_URL, $nltk_py_path);
        curl_exec($curl);
    
        //$handle = popen($cmdline, "r");

        //$output = fread($handle, 4096);

        //print_r($output);

        curl_close($curl);

        $this->updateQuizStatus($quiz->keyword, $quiz->source, STATUS::QG_COMPLETED);

        return $quiz;
    }

    public function getAllQuizzes($start, $end, $sort)
    {
        $quizzes = array();

        $order = " keyword ";

        switch($sort) {
            default:
                $sort = 'alpha';
                break;
            case "new":
                $order = 'when';
                break;
            case "most_read":
                $order = 'total_read_count';
                break;
            /*case "trending":
                //calculate from request table with timestamp within last 'x' months
                $
                //$order = " "
                break;
            */
        }

        $sql = "SELECT * FROM quizzes WHERE status=:status ORDER BY :order LIMIT :start, :row_count";

        //print_r($sql);

        $count_sql = "SELECT COUNT(*) FROM quizzes";

        try {

            //mysql doesnt allow subquery as limit clause, and this allows us to specify a row_count s.t.
            //start+row_count < end
            $stmt = Connection::$conn->prepare($count_sql);
            $stmt->execute();
            $rows = $stmt->fetchAll(PDO::FETCH_ASSOC);

            $num_rows = (int)$rows[0]["COUNT(*)"];
            #print_r($num_rows);
            if($end == -1 || $num_rows < (int)$end) {
                $end = $num_rows;
            }

            $stmt = Connection::$conn->prepare($sql);

            if(isset($start)) {
                $stmt->bindValue(":start", (int)$start, PDO::PARAM_INT);
            } else {
                $stmt->bindValue(":start", 0, PDO::PARAM_INT);
            }

            $row_count = $end-$start;

            $stmt->bindValue(":row_count", $row_count, PDO::PARAM_INT);
            $stmt->bindValue(":order", $order, PDO::PARAM_STR);
            $stmt->bindValue(":status", STATUS::READY, PDO::PARAM_INT);

            $stmt->execute();

            $rows = $stmt->fetchAll(PDO::FETCH_ASSOC);

            foreach($rows as $row) {
                $quizzes[] = new Quiz($row);
            }

        } catch (Exception $e) {

            echo "Error getting questions: " + $e->getMessage();
        }

        return array("sort"=> $sort, "quizzes" => $quizzes, "start"=>$start, "end"=>$end);
    }

    public function getRandomQuizzes($numberQuizzes)
    {

    if(!isset(Connection::$conn)) {
        echo "connection null";
    }

    $quizzes = array();

        $sql = "SELECT * FROM quizzes ORDER BY RAND() LIMIT :numberQuizzes";

        try {

            $stmt = Connection::$conn->prepare($sql);

            $stmt->bindValue(":numberQuizzes", $numberQuizzes, PDO::PARAM_INT);


            $stmt->execute();

            $rows = $stmt->fetchAll(PDO::FETCH_ASSOC);


            foreach($rows as $row) {
                $quizzes[] = new Quiz($row);
            }

        } catch (Exception $e) {

            echo "Error getting questions: " + $e->getMessage();
        }

        return $quizzes;
    }

    public function getRandomAnswers($numberAnswers)
    {
        $answers = array();

        $sql = "SELECT * FROM quiz_questions WHERE q_type=1 AND CHAR_LENGTH(answer)>1 ORDER BY RAND() LIMIT :numberAnswers";

        try {

            $stmt = Connection::$conn->prepare($sql);

            $stmt->bindValue(":numberAnswers", $numberAnswers, PDO::PARAM_INT);


            $stmt->execute();

            $rows = $stmt->fetchAll(PDO::FETCH_ASSOC);


            foreach($rows as $row) {
                $answers[] = $row['answer'];
            }

        } catch (Exception $e) {

            echo "Error getting questions: " + $e->getMessage();
        }

        return $answers;
    }

    public function getQuizQuestions($keyword, $source)//, $lang)
    {
        $quiz = $this->findQuiz($keyword, $source);//, $lang);

        #print_r($quiz);
        $sql = "SELECT * FROM quiz_questions WHERE quiz_id=:quiz_id LIMIT 10";

        try {

            $stmt = Connection::$conn->prepare($sql);

            $stmt->bindValue(":quiz_id", $quiz->id);

            $stmt->execute();

            $results = $stmt->fetchAll(PDO::FETCH_ASSOC);


            $questions = array();
            foreach($results as $result) {

                do {
                    $otherAnswers = $this->getRandomAnswers(3);
                } while(in_array($result["answer"], $otherAnswers));

                $chr = mb_substr($result["answer"], 0, 1, "UTF-8");
                $capitalize = mb_strtolower($chr, "UTF-8") != $chr;

                $case_fn = "mb_strtolower";
                if($capitalize) {
                    $case_fn = "mb_strtoupper";
                }

                    for( $i = 0; $i < sizeof($otherAnswers); $i++) {
                        #print_r($otherAnswers[$i]);
                        $otherAnswers[$i] = call_user_func($case_fn, mb_substr($otherAnswers[$i], 0, 1, "UTF-8")) . mb_substr($otherAnswers[$i], 1, mb_strlen($otherAnswers[$i]), "UTF-8");
                        #print_r($otherAnswers[$i]);
                    }

                //preg_replace('\\n', '',
                // Delimiter must not be alphanumeric or backslash[
                $answers = array(
                    $result["answer"],
                    $otherAnswers[0],
                    $otherAnswers[1],
                    $otherAnswers[2],
                );

                shuffle($answers); //DevSkim: ignore DS148264 

                //preg_replace('\\n', '',
                // Delimiter must not be alphanumeric or backslash[
                $questions[] = array(
                    "question" 	=> $result["question"],
                    "q_type"    => $result["q_type"],
                    "id"        => $result["id"],
                    "answer" 	=> ($result["q_type"]==1)?$answers:[],
                );

            }


            shuffle($questions); //DevSkim: ignore DS148264 

            $quiz->questions = $questions;

	#	print_r($quiz);

        } catch (Exception $e){
            echo "Error reading or parsing existing quiz questions: " . $e->getMessage();
        }


		return $quiz;
	}

    public function getStudyQuestions($keyword, $source)//, $lang)
    {
        $quiz = $this->findQuiz($keyword, $source);//, $lang);

        #print_r($quiz);
	    $sql = "SELECT * FROM quiz_questions WHERE quiz_id=:quiz_id LIMIT 10";

        try {

            $stmt = Connection::$conn->prepare($sql);

            $stmt->bindValue(":quiz_id", $quiz->id);

            $stmt->execute();

            $results = $stmt->fetchAll(PDO::FETCH_ASSOC);


            $questions = array();
            foreach($results as $result) {

                do {
                    $otherAnswers = $this->getRandomAnswers(3);
                } while(in_array($result["answer"], $otherAnswers));

                $chr = mb_substr($result["answer"], 0, 1, "UTF-8");
                $capitalize = mb_strtolower($chr, "UTF-8") != $chr;

                $case_fn = "mb_strtolower";
                if($capitalize) {
                    $case_fn = "mb_strtoupper";
                }

                    for( $i = 0; $i < sizeof($otherAnswers); $i++) {
                        #print_r($otherAnswers[$i]);
                        $otherAnswers[$i] = call_user_func($case_fn, mb_substr($otherAnswers[$i], 0, 1, "UTF-8")) . mb_substr($otherAnswers[$i], 1, mb_strlen($otherAnswers[$i]), "UTF-8");
                        #print_r($otherAnswers[$i]);
                    }

                //preg_replace('\\n', '',
                // Delimiter must not be alphanumeric or backslash[
                $answers = array(
                    $result["answer"],
                    $otherAnswers[0],
                    $otherAnswers[1],
                    $otherAnswers[2],
                );

                shuffle($answers); //DevSkim: ignore DS148264 

                //preg_replace('\\n', '',
                // Delimiter must not be alphanumeric or backslash[
                $questions[] = array(
                    "question" 	=> $result["question"],
                    "q_type"    => $result["q_type"],
                    "id"        => $result["id"],
                    "answer" 	=> ($result["q_type"]==1)?$answers:[],
                    "p_id"      => $result["p_id"],
                );

            }


            shuffle($questions); //DevSkim: ignore DS148264 

            $quiz->questions = $questions;

	#	print_r($quiz);

        } catch (Exception $e){
            echo "Error reading or parsing existing quiz questions: " . $e->getMessage();
        }


		return $quiz;
	}

	public function grade($user, $questionID, $answer)
    {
        $boolResult = false;

        $sql = "SELECT * from quiz_questions WHERE id=:id;";

        try {

            $stmt = Connection::$conn->prepare($sql);

            $stmt->bindValue(":id", $questionID);

            $stmt->execute();

            $row = $stmt->fetch(PDO::FETCH_ASSOC);

            $boolResult = (!strcmp($row['answer'], $answer))?true:false;

        } catch(Exception $e){
            
            echo json_encode(array("status" => "error grading: " . $e->getMessage()));
        }
/*
        $sql = "select * from quiz_users WHERE username=:username;";

        try {

            $stmt = Connection::$conn->prepare($sql);

            $stmt->bindValue(":username", $user);

            $stmt->execute();
            $row = $stmt->fetch(PDO::FETCH_ASSOC);

            $streak = intval($row['streak']);
            $best_streak = intval($row['best_streak']);
            $totalscore = intval($row['totalscore']);

            if($boolResult)
            {
                $streak++;

                if($streak > $best_streak)
                    $best_streak++;

                $totalscore++;
            } else {
                $streak = 0;
            };

            $sql = "UPDATE quiz_users set best_streak=:best_streak, streak=:streak, totalscore=:totalscore where username=:username;";

            $stmt = Connection::$conn->prepare($sql);

            $stmt->bindValue(":best_streak", $best_streak, PDO::PARAM_INT);
            $stmt->bindValue(":streak", $streak, PDO::PARAM_INT);
            $stmt->bindValue(":totalscore", $totalscore, PDO::PARAM_INT);

            $stmt->bindValue(":username", $user);

            $stmt->execute();

        } catch(Exception $e){

            echo json_encode(array("status" => "error grading: " . $e->getMessage()));

        }
*/
        return $boolResult;
    }
}
?>
