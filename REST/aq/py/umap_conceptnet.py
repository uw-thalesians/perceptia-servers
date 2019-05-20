#!/usr/bin/env python2.7

print "Content-Type: application/json\n\n"

import sys
import requests as r
import cgi, cgitb, json, hashlib
import umap
import mysql_quiz
#from nltk import wordnet
#from pyld import jsonld
import pandas as pd

cgitb.enable()

data = cgi.FieldStorage()

vals = {}

def convertCGI(data):
  for key in data.keys():
    vals[key] = []

    for list_val in data.getlist(key):
      vals[key].append( list_val )

convertCGI(data)

if "root" in vals:
  print vals["root"]


quiz_keywords = mysql_quiz.get_all_keywords()

num_nodes = len(quiz_keywords)
if num_nodes==0:
  print json.dumps({"umap": {}, "num_nodes" : 0})
  quit()

provided_n_neighbors = 20

if "n_neighbors" in vals:
  provided_n_neighbors = vals["n_neighbors"]

if num_nodes < provided_n_neighbors:
  provided_n_neighbors = num_nodes

provided_min_dist = 0.1
if "min_dist" in vals:
  provided_min_dist = vals["min_dist"]


#https://github.com/commonsense/conceptnet5/wiki/API
#http://api.conceptnet.io/c/en/

#print quiz_keywords

concepts = []

for keyword in quiz_keywords:
    #print keyword[0]
    url_component = cgi.escape(keyword.replace(" ","_"))
    resp_json = r.get("http://api.conceptnet.io/c/en/"+url_component).json()
    
    if "error" in resp_json:
      concepts.append({"keyword":keyword, "error": resp_json["error"]["status"]})


    en_only = [val for val in resp_json["edges"]\
                   if "language" in val["start"] and\
                       val["start"]["language"]=="en" and\
                       "language" in val["end"] and\
                       val["end"]["language"]=="en"]

    for entry in en_only:
        #in the sense of meronym, for /has a/ et c. relationships
        sense = "part"
        #self-reference for uses of, e.g. coffee is a drink and a can be used to refer to a color
        part = "self"
        
        if "end" in entry:
            if "sense_label" in entry["end"]:
                sense = entry["end"]["sense_label"]

            if "label" in entry["end"]:
                part = entry["end"]["label"]
                
            rel = entry["rel"]["@id"]

        concepts.append({"keyword":keyword, "rps":rel+part+sense, "sense": sense})#,  "part": part,  , "relationship": rel})

concept_df = pd.DataFrame(concepts)
#print concept_df
#if only errors exist, just respond with those errors
if "error" in concept_df.columns and concept_df.drop("error", axis=1).shape[1]==1:
  print json.dumps({"umap":[{"keyword": row["keyword"], "error": row["error"]} for row in concept_df.to_dict(orient="index").values()]})
  quit()

li_vectors = pd.get_dummies(concept_df["sense"])
li_vectors["keyword"] = concept_df["keyword"]
vectors = li_vectors.groupby("keyword").sum()

init_algo = "spectral"
if num_nodes < len(vectors.columns):
  init_algo = "random"

reducer = umap.UMAP(random_state=0, min_dist=provided_min_dist, n_neighbors=provided_n_neighbors, init=init_algo)

embedding = reducer.fit_transform(vectors.values)

vectors["dom_sense"] = vectors.idxmax(axis=1)

def hash_to_color(id):
    h = hashlib.md5(id.encode("UTF-8"))
    return "#"+h.hexdigest()[0:18:3]+"FF"

color_dict = {}
for sense in vectors["dom_sense"].unique():
    color_dict[sense] = hash_to_color(sense)

vectors["color"] = vectors["dom_sense"].apply(lambda x: color_dict[x])


vectors["x"] = embedding[:,0]
vectors["y"] = embedding[:,1]

print json.dumps({"umap": vectors[["x", "y", "color"]].to_dict(orient="index"), "num_nodes":num_nodes})