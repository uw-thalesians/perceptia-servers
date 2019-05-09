docker kill aqsolr
docker rm aqsolr
docker volume rm solr
docker run -d -v solr:/var/solr -p 8983:8983 --network perceptia-net --name aqsolr uw-thalesians/solrstream solr-precreate stream /opt/solr/server/solr/configsets/stream