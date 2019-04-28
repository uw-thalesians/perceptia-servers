<?php

require_once 'config.php';

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

    public function findQuiz($keyword)
    {

        $sql = "SELECT * FROM quizzes WHERE keyword=:keyword";

        $quiz = null;

        try {

            $stmt = Connection::$conn->prepare($sql);

            $stmt->bindValue(":keyword", $keyword, PDO::PARAM_STR);

            $stmt->execute();

            $row = $stmt->fetch(PDO::FETCH_ASSOC);

            if($row) {

                $quiz = new Quiz( $row );

            } else {

                $summary = $this->fetchSummary($keyword);

                $quiz = $this->addNewQuiz($keyword, $summary);
            }

        } catch (Exception $e) {
            echo "Error finding quiz: " . $e->getMessage();
        }

        return $quiz;
    }

    private function fetchSummary($keyword) {

        //curl call py script
        $curl_handle= curl_init();
        $server = "localhost";
        $wikipedia_python_rest_query = $server.$this->_f3->get("BASE")."/py/w_rest.py?keyword=" . urlencode($keyword);
        #$wikipedia_python_rest_query = $server . "/?keyword=" . urlencode($keyword);
        #print_r($wikipedia_python_rest_query);

        curl_setopt($curl_handle, CURLOPT_URL, $wikipedia_python_rest_query);
        curl_setopt($curl_handle, CURLOPT_RETURNTRANSFER, 1);
        #curl_setopt($curl_handle, CURLOPT_PORT, 27277);

        $json = curl_exec($curl_handle);
        #print_r(curl_error($curl_handle));
        #print_r(curl_getinfo($curl_handle, CURLINFO_HTTP_CODE));
        curl_close($curl_handle);
        #print_r($json);

        $summary = json_decode($json, true)["summary_text"];

        #print_r($summary);

        return $summary;
    }

    private function addNewQuiz($keyword, $summary)//, $lang)
    {
        $sql = "INSERT INTO quizzes (keyword, image, summary) VALUES (:keyword, :image, :summary)";

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
        $google_cse_rest_api_get = "https://www.googleapis.com/customsearch/v1?q=%22$search_keyword%22&cx=004799634748936919555:ewzgppgp6wu&searchType=image&safe=high&key=AIzaSyAjG8-dIi5U_SeOKoaTMuHYseILeSSmRgU";

        //print_r($google_cse_rest_api_get);

        //Google Custom Search v1 REST API
        ////curl_setopt($curl, CURLOPT_URL, $google_cse_rest_api_get);
        ////curl_setopt($curl, CURLOPT_RETURNTRANSFER, 1);

        ////$json = curl_exec($curl);

        //print_r($json);

        ////$search = json_decode($json, true);
        $search = array();
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

                    curl_setopt($curl, CURLOPT_HEADER, 0);
                    curl_setopt($curl, CURLOPT_BINARYTRANSFER, 1);
                    curl_setopt($curl, CURLOPT_URL, $imageURL);

                    curl_setopt($curl, CURLOPT_FOLLOWLOCATION, true);

                    $localPath = 'images/' . $filename;


                    if( file_exists($localPath) ) {
                        //try to create a unique filename
                        $pathComponents = pathinfo($localPath);

                        $localPath = 'images/' . $pathComponents['filename'] . time() . "." . $pathComponents['extension'];

                        //print_r($localPath);

                        //if somehow something still went wrong, use default quiz.png
                        if(file_exists($localPath))
                            unset($localPath);
                    }

                    //var_dump($localPath);

                    if(isset($localPath)) {
                        $localFile = fopen($localPath, 'w');

                        curl_setopt($curl, CURLOPT_FILE, $localFile);
                        curl_exec($curl);
                        fclose($localFile);
                    }
                }
                $tries++;
            } while((filesize($localPath) < 10000) && $tries < 10);

            curl_close($curl);
        }else{
            //var_dump($json);
        }

        if(!isset($localPath))
            $localPath = 'images/quiz.png';

        try {

            $stmt = Connection::$conn->prepare($sql);

            $stmt->bindValue(":keyword", $keyword);
            $stmt->bindValue(":summary", $summary);
            $stmt->bindValue(":image", $localPath);
            //$stmt->bindValue(":lang", $lang);

            $stmt->execute();

            $quiz = $this->findQuiz($keyword);//, $lang);

        } catch (Exception $e){
            echo "Error creating quiz: " . $e->getMessage();
        }

        $path_to_py = dirname($_SERVER['PHP_SELF']) . "/py/n_cgi.py";

        #print_r($path_to_py);

        $nltk_py_path = "http://${_SERVER['SERVER_NAME']}${path_to_py}?keyword=" . urlencode($quiz->keyword);// . "&lang=" . urlencode($lang);

        #print_r($nltk_py_path);
        $curl = curl_init();
        curl_setopt($curl, CURLOPT_URL, $nltk_py_path);
        curl_exec($curl);
    
        //$handle = popen($cmdline, "r");

        //$output = fread($handle, 4096);

        //print_r($output);

        curl_close($curl);

        return $quiz;
    }

    public function getAllQuizzes()
    {
        $quizzes = array();

        $sql = "SELECT * FROM quizzes ORDER BY keyword";

        try {

            $stmt = Connection::$conn->prepare($sql);

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

    public function getQuizQuestions($keyword)//, $lang)
    {
        $quiz = $this->findQuiz($keyword);//, $lang);

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

            //header("Content-type: application/json");
            //echo json_encode(array("answer"=>$answer, "realanswer"=>$row['answer'], "result"=>$boolResult));

        } catch(Exception $e){
            header("Content-type: application/json");
            echo json_encode(array("status" => "error grading: " . $e->getMessage()));
        }

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

            header("Content-type: application/json");
            echo json_encode(array("status" => "error grading: " . $e->getMessage()));

        }

        return $boolResult;
    }
}
?>
