#!/usr/bin/env python2.7
print "Content-Type: application/json\n\n"

import sys
import cgi, cgitb, json
cgitb.enable()

data = cgi.FieldStorage()

vals = {}

def convertCGI(data):
  for key in data.keys():
    vals[key] = []

    for list_val in data.getlist(key):
      vals[key].append( list_val )

convertCGI(data)

#print json.loads(json.dumps(vals))
import n

if "keyword" not in vals:
  print "{\"error\":\"Parsing of keyword GET parameter failed\"}"
  exit

keyword_str = " ".join(vals["keyword"])

n.create_quiz(keyword_str)