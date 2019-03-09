FROM mcr.microsoft.com/mssql/server:2017-CU12-ubuntu

# Create directory to hold db scripts
RUN mkdir -p /script
WORKDIR /script

# Copy db files
COPY Perceptia/ /script
COPY setup-db.sh /script
COPY entrypoint.sh /script

# Make setup script executable
RUN chmod +x setup-db.sh
RUN chmod +x entrypoint.sh

CMD /bin/bash ./entrypoint.sh

