#!/usr/bin/env python2.7

import sys, os, random, math, json

#path = os.path.dirname(__file__)+"/user-site-packages"

#print(path)
#sys.path.append("./lib/python2.7/site-packages")
try:
    import nltk, mysql_quiz, re
    
    from nltk.corpus import wordnet#, conll2002
    import spacy
    #with 2G mem and 2G swap, md can be loaded
    nlp = spacy.load("en_core_web_md")
    #with 2.5G mem and 2.5G swap, vectors_web_lg can be loaded
    #this by itself throws a sentencizer/parser error? sentence boundaries? is this to be used with one of the others?
    #nlp = spacy.load("en_vectors_web_lg")
except Exception as e:
    print json.dumps({"error":str(e)})
    sys.exit(2)

nltk.data.path.append("./nltk_data")

try:
    #import summarize#, w2v
    pass
except Exception as e:
    print json.dumps({"error":str(e)})

def find_alternative_responses(sentence, correct_answer, num_resp=1, keyword=None):
    responses = []

    #try to find as many as are asked for, but if not found don't get stuck looking
    for cur_resp_num in range(0, num_resp):
        #default to true answer
        replacement = answer[-1][1]

        found_replacement = False
        #if false was selected, find an alternative answer
        if chosen_answer == "f":
            #get the first wordnet synset for this word
            synsets = wordnet.synsets(answer[-1][1], pos="n")

            print("looking for a replacement for {}".format(answer[-1][1]))

            if len(synsets) > 0:
                for synset in synsets:
                    print("actual answer sense:{}\n".format(synset))
                    for hypernym in synset.hypernyms():
                        #try using keyword similarity to find more likely answers by related hyponyms?
                        #print("sense hypernym: {}\n".format(hypernym))

                        hyponyms = hypernym.hyponyms()
                        #print("hyponyms of this hypernym: {}\n".format(hyponyms))
                        #for hyponym in hyponyms:
                        try:
                            hyponyms.remove(synset)
                        except Exception as e:
                            pass
                        if len(hyponyms)>0:
                            replacement_syn = hyponyms[int(math.floor(random.random()*len(hyponyms)))]
                            #if hyponym == synset:
                            #    print("skipping same synset hyponym")
                            #    continue
                            if len(replacement_syn.lemmas()) > 0:
                            #for lemma in replacement_syn.lemmas():
                                lemma = replacement_syn.lemmas()[int(math.floor(random.random()*len(replacement_syn.lemmas())))]
                                replacement = lemma.name()
                                print("setting replacement {}".format(replacement))
                                #if (len(replacement) > len(orig_answer)-3) or (len(replacement) < len(orig_answer)+3):
                                #    break
                                found_replacement = True
                                print("found_replacement {}".format(found_replacement))

            #if there are no good alternatives we found, just make this a true question instead of having no additional way to handle finding a way to make it false (maybe add negation later? e.g. inject "not "+answer[-1][1]
        
    return responses

def create_question(sentence, keyword=None):

    question = []

    #print(u"sentence: {} | keyword: {}".format(sentence, keyword))

    if len(sentence.text) < 3:
        return u"", u"", 0

    q_type = int(math.floor(random.random()*2)+1) #should this use mysql_quiz to get types num types and then match to a generation scheme?
    

    #print("about to get nouns")

    question, answer = get_all_nouns(sentence.text, keyword)

    #print(u"question: {} | answer: {}".format(question, answer))
    #print("got nouns")

    #if no suitable word was found, just move on to next possible question
    if len(answer) < 1:
        return u"", u"", 0

    #question = [word[0] for word in question]

    #print question

    answer_index = int(math.floor(random.random()*len(answer)))

    #print(answer)

    word_index_start = answer[answer_index][0]
    word_index_end = answer[answer_index][1]
    chosen_answer = answer[answer_index][2]


    if q_type == 1:
    #print answer, word_index

        question = question[:word_index_start] + [u"__________"] + question[word_index_end:]

        #print(question)
        #answer = answer[-1][1]

    else:
        #print("making tf q")
        chosen_answer = "f" if math.floor(random.random()*2)==0 else "t"

        orig_answer = answer[answer_index][2]
        #default to true answer
        replacement = answer[answer_index][2]

        found_replacement = False
        #if false was selected, find an alternative answer
        if chosen_answer == "f":
            #get the first wordnet synset for this word
            synsets = wordnet.synsets(answer[answer_index][2], pos="n")

            #print("looking for a replacement for {}".format(answer[answer_index][2]))

            if len(synsets) > 0:
                for synset in synsets:
                    #print("actual answer sense:{}\n".format(synset))
                    for hypernym in synset.hypernyms():
                        #try using keyword similarity to find more likely answers by related hyponyms?
                        #print("sense hypernym: {}\n".format(hypernym))

                        hyponyms = hypernym.hyponyms()
                        #print("hyponyms of this hypernym: {}\n".format(hyponyms))
                        #for hyponym in hyponyms:
                        try:
                            hyponyms.remove(synset)
                        except Exception as e:
                            pass
                        if len(hyponyms)>0:
                            replacement_syn = hyponyms[int(math.floor(random.random()*len(hyponyms)))]
                            #if hyponym == synset:
                            #    print("skipping same synset hyponym")
                            #    continue
                            if len(replacement_syn.lemmas()) > 0:
                            #for lemma in replacement_syn.lemmas():
                                lemma = replacement_syn.lemmas()[int(math.floor(random.random()*len(replacement_syn.lemmas())))]
                                replacement = lemma.name().replace("_", " ")
                                #print("setting replacement {}".format(replacement))
                                #if (len(replacement) > len(orig_answer)-3) or (len(replacement) < len(orig_answer)+3):
                                #    break
                                found_replacement = True
                                #print("found_replacement {}".format(found_replacement))

            #if there are no good alternatives we found, just make this a true question instead of having no additional way to handle finding a way to make it false (maybe add negation later? e.g. inject "not "+answer[-1][1]
            #print("chosen answer {} and found_replacement {}".format(chosen_answer, found_replacement))
            if (chosen_answer == "f") and (found_replacement == False):
                #print("a replacement wasn't found, converting to true")
                chosen_answer = "t"
                #replacement = 

        #print(u"question to rebuild for tf: {}".format(question))
        question = question[:word_index_start] + [replacement] + question[word_index_end:]

    question = u" ".join(question)

    return question, chosen_answer, q_type

def filter_noun_chunks(doc):
    better_noun_chunks = []
    for chunk in doc.noun_chunks:
        usable_chunk = chunk
        good = True
        for token in chunk:
            #print(token.text, token.pos_)
            #remove chunks that use pronouns
            if token.pos_ == "PRON":
                #print(u"skipping inclusion of {}".format(chunk.text))
                good = False
                break
            #if it begins with a determiner, deepcopy and update to remove
            elif token.pos_ == "DET":
                #print(u"found determiner")
                if chunk.start+1 < chunk.end:
                    usable_chunk = doc[chunk.start+1:chunk.end]
                    #print(u"using {}".format(usable_chunk.text))
                else:
                    good = False
                break
        if good:
            better_noun_chunks.append(usable_chunk)

    return better_noun_chunks

def get_all_nouns(sentence, keyword):
    doc = nlp(sentence)
    question = []
    try:
        better_noun_chunks = filter_noun_chunks(doc)

        answer = [[chunk.start, chunk.end, chunk.text] for chunk in better_noun_chunks if keyword.lower() not in chunk.text.lower()]

        question = [token.text for token in doc]#"and now for something completely different")
        #question = 
    except Exception as e:
        print json.dumps({"error":str(e)})
    
    #answer = [[i, word[0]] for i, word in enumerate(question) if 'NN' in word[1]]
    return question, answer

#def get_noun(sentence):
#    print sentence
#    noun = get_all_nouns(sentence)
#    print noun
#    if len(noun) > 0:
#        noun=noun[1][-1][1]
#
#
#    return noun

def get_sentences(text):
    doc = nlp(text)
    return doc.sents

def create_quiz(keyword):
    
    summary, quiz_id = mysql_quiz.get_text(keyword)
    
    for p_id, paras in summary.items():
        #for para in paras:
            #print "para_id", p_id, "all_paras", paras
            usum = paras.decode("UTF-8")

            sentences = get_sentences(usum)

            for sentence in sentences:
                question, answer, q_type = create_question(sentence, keyword)

                if question == u"" or answer == u"":
                    continue

                mysql_quiz.add_question( question, answer, quiz_id, q_type, p_id )


def create_graph_quiz(title):
    (nodes, graph_quiz_id) = mysql_quiz.get_graph_quiz(title)

    for node in nodes:
        (summary, quiz_id) = mysql_quiz.get_text(keyword)

        sentences = summarize.summarize_text(unicode(summary)).summaries

        w2v.learn_new_sentences(sentences)

        for sentence in sentences:

            question, answer = create_question(sentence)

            if question == "" or answer == "":
                continue

            mysql_quiz.add_question( question, answer, quiz_id )
