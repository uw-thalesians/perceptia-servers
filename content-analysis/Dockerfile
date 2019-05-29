FROM solr:8
SHELL ["/bin/bash", "-c"]
USER root
RUN apt update && apt install -y less vim net-tools
USER solr
COPY stream /opt/solr/server/solr/configsets/stream/
RUN ls /opt/solr/server/solr/configsets
