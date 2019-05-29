FROM mcr.microsoft.com/mssql/server:2017-CU12-ubuntu

# Create directory to hold db scripts
RUN mkdir -p /script
RUN mkdir -p /script/Perceptia
WORKDIR /script

# Copy db files
COPY Perceptia/ /script/Perceptia
COPY ./script/setup-db.sh /script
COPY ./script/entrypoint.sh /script

# Make setup script executable
RUN chmod +x setup-db.sh
RUN chmod +x entrypoint.sh

CMD /bin/bash /script/entrypoint.sh

