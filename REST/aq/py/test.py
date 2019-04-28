#!/usr/bin/env python2.7
print "Content-Type: text/html\n\n"
import cgi, cgitb, json, os
cgitb.enable()

data = cgi.FieldStorage()

vals = {}

def convertCGI(data):
  for key in data.keys():
    vals[key] = []
    #print key, data.getlist(key)
    for list_val in data.getlist(key):
      #print list_val
      vals[key].append( list_val )
      #convertCGI(data[key])
    
    #vals[key] = data[key].value

convertCGI(data)

print json.loads(json.dumps(vals))
#print dict([data.getlist(key)[i] for key in data.keys() for i in range(len(key)-1) ] )#json.encode(data)

#import sys
#sys.stderr = sys.stdout
#print "test"
#cgi.test()
print(os.environ.get("HTTP_X_FORWARDED_FOR"))
