#!/usr/bin/env python2.7

print "Content-Type: text/html; charset=UTF-8\n\n"

import sys, os, json, cgi, cgitb

try:
  user_path = os.getcwd()+"/lib/python2.7/site-packages"
  sys.path.append(user_path)
  import wikipedia
  from wikipedia import DisambiguationError
  from bs4 import BeautifulSoup
  #from langdetect import detect
  import codecs
  #import mysql_quiz
  #import unicodedata
except Exception as e:
  print unicode(e)
  sys.exit()

data = cgi.FieldStorage()

keyword = data["keyword"].value
#lang = data.getvalue("lang", "en")

#text = keyword
try:
  pass
  #utf8reader = codecs.getreader("utf8")
  #keyword = utf8reader(codecs.StreamReader(keyword)).readline()
  #print(keyword)
  #lang = detect(keyword)
  #print(lang)
except Exception as e:
  print(unicode(e))
  pass

try:

  #wikipedia.set_lang(lang)

  keyword = wikipedia.search(keyword)
  if len(keyword)>0:
    keyword=keyword[0]
 
  try:
    page = wikipedia.page(keyword)
  except DisambiguationError as e:
    #print(e.options)
    keyword = e.options[0]
    page = wikipedia.page(keyword)

  #print summary_text
  #summary_text_short = wikipedia.summary(keyword)
  #print("about to use beautifulsoup")

  bs4_page = BeautifulSoup(page.html(), "html.parser")

    

  [sup.extract() for sup in bs4_page("sup", attrs={"class":"reference"})]
  para = [p.text for p in bs4_page.select("p")]
  #print("processed")
  summary_text = para
  #summary_text = u" ".join(para)
  #print("selected para")
  #print(summary_text)

  #summary_text = summary_text.replace("\n", " ")
  #summary_text = full_text.replace("\n", " ")
  #summary_text = unicodedata.normalize("NFKD", summary_text).encode('ascii', 'ignore')
  #print summary_text
except Exception as e:
  #summary_text = u" ".join(para)
  summary_text = unicode(e)

#print summary_text
#print json.dumps({"summary_text":summary_text})
print json.dumps({u"summary_text":summary_text, u"requested_from":os.environ.get("HTTP_X_FORWARDED_FOR")})
#mysql_quiz.put_text(keyword, summary_text)
