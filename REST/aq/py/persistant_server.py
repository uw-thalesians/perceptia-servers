#!/usr/bin/env python2.7
import SimpleHTTPServer
import SocketServer

import sys, os, json

try:
  user_path = os.getcwd()+"/lib/python2.7/site-packages"
  sys.path.append(user_path)
  import wikipedia
  import spacy
  nlp = spacy.load("en_core_web_sm")
  from wikipedia import DisambiguationError
  from bs4 import BeautifulSoup
  #from langdetect import detect
  #import codecs
  #import mysql_quiz
  #import unicodedata
except Exception as e:
  print unicode(e)
  sys.exit()

class WikiPageRequestHandler(SimpleHTTPServer.SimpleHTTPRequestHandler):
    
    def do_GET(self):
        kv = {}
        key_val_strings = self.path[self.path.index("?")+1:].split("&")
        for pair in key_val_strings:
            split_pair = pair.split("=")
            kv[split_pair[0]] = split_pair[1]

        page = wikipedia.page(kv["keyword"])

        self.protocol_version = 'HTTP/1.1'
        self.send_response(200, 'OK')
        self.send_header('Content-type', 'application/json')
        self.end_headers()
        self.wfile.write(bytes(json.dumps({u"summary_text":page.content})))
        return

PORT = 27277

Handler = WikiPageRequestHandler

httpd = SocketServer.TCPServer(("", PORT), Handler)

httpd.server_name = ""
httpd.server_port = PORT

print "serving at port", PORT
httpd.serve_forever()
