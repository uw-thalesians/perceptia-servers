#!/usr/bin/env python2.7

import json
import mysql.connector

from config import mysql_user, mysql_pass, mysql_host, mysql_db, mysql_port

def get_all_text():

    try:
        cnx = mysql.connector.connect(user=mysql_user, password=mysql_pass, host=mysql_host, port=mysql_port, database=mysql_db)
        curRetrieveQuiz = cnx.cursor(buffered=True)

        query = "select * from paragraphs"

        curRetrieveQuiz.execute(query)

        summaries = []

        for (curSummary) in curRetrieveQuiz:
            summaries.append(curSummary)

    except mysql.connector.Error as e:
        print json.dumps({"error":str(e)})
    finally:
        cnx.close()

    return summaries

def get_all_keywords():
    try:
        cnx = mysql.connector.connect(user=mysql_user, password=mysql_pass, host=mysql_host, port=mysql_port, database=mysql_db)
        curRetrieveKeyword = cnx.cursor(buffered=True)

        query = "select keyword from quizzes"

        curRetrieveKeyword.execute(query)

        keywords = []

        for (curKeyword) in curRetrieveKeyword:
            keywords.append(curKeyword)

    except mysql.connector.Error as e:
        print(str(e))
    finally:
        cnx.close()

    return keywords

def put_text(keyword, summary):
    row_id=None
    source = 'wiki'
    try:
        cnx = mysql.connector.connect(user=mysql_user, password=mysql_pass, host=mysql_host, port=mysql_port, database=mysql_db)
        cursor = cnx.cursor()

        #print("about to insert from python")

        query = "insert into quizzes (keyword, source, total_read_count) VALUES (%s, %s, %s, %s)"
        vals = (keyword)

        cursor.execute(query, vals, source, 1)

        row_id = cursor.lastrowid
        #print(curRetrieveQuiz)
        #for (curKeyword, curSummary, curQuiz_id) in curRetrieveQuiz:
        #    print(curKeyword, curSummary, curQuiz_id)
            #keyword = curKeyword
        #    summary = curSummary
        #    quiz_id = curQuiz_id

    except mysql.connector.Error as e:
        print(json.dumps({"error":str(e)}))
    finally:
        cnx.close()

    return row_id

def get_text(keyword):#, lang):
    #print "get_text", keyword
    summary = ""
    quiz_id = ""

    try:
        cnx = mysql.connector.connect(user=mysql_user, password=mysql_pass, host=mysql_host, port=mysql_port, database=mysql_db)
        curRetrieveQuiz = cnx.cursor(buffered=True)

        query = "select keyword, id from quizzes where keyword='"+keyword+"'"# AND lang='"+lang+"'"

        curRetrieveQuiz.execute(query)

        curRetrieveParas = cnx.cursor(buffered=True)

        for (curKeyword, curQuiz_id) in curRetrieveQuiz:

            query = "select id, text from paragraphs where quiz_id='"+str(curQuiz_id)+"'"# AND lang='"+lang+"'"

            curRetrieveParas.execute(query)

            summary = {}

            for (p_id, text) in curRetrieveParas:

                summary[p_id] =text

            quiz_id = curQuiz_id

    except mysql.connector.Error as e:
        print json.dumps({"error":str(e)})
    finally:
        cnx.close()

    return summary, quiz_id

def get_graph_quiz(title):
    
    try:
        cnx = mysql.connector.connect(user=mysql_user, password=mysql_pass, host=mysql_host, port=mysql_port, database=mysql_db)
    
	curRetrieveQuiz = cnx.cursor(buffered=True)

        query = "select id from graph_quizzes where title='"+title+"'"

        curRetrieveQuiz.execute(query)

        graph_quiz_id = curRetrieveQuiz[0][0]

        query = "select id, keyword from graph_nodes where graph_quiz_id='"+graph_quiz_id+"'"

	curRetrieveQuiz.execute(query)
        nodes = {}

        for (id, keyword) in curRetrieveQuiz:
             nodes[keyword]=id

    except Exception as e:
        print(str(e))
    finally:
        cnx.close()
    return nodes, graph_quiz_id

def add_question(question, answer, quiz_id, q_type, p_id):
    row_id = None
    try:
        cnx = mysql.connector.connect(user=mysql_user, password=mysql_pass, host=mysql_host, port=mysql_port, database=mysql_db)

        curCreateQuestion = cnx.cursor()#buffered=True)

        curCreateQuestion.execute("insert into quiz_questions (question, answer, quiz_id, q_type, p_id) values (%s, %s, %s, %s, %s)", (question, answer, quiz_id, q_type, p_id))

        cnx.commit()
        
        row_id = curCreateQuestion.lastrowid
    except mysql.connector.Error as e:
            print(e)
    finally:
        cnx.close()

    return row_id
