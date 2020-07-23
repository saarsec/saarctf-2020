# Build things
# You can change the base image here, but you also need to change it in .gitlab-ci.yml
FROM node:lts-buster AS buildstage
ADD ./service /opt/service
ADD ./build.sh /opt/
WORKDIR /opt/
RUN ./build.sh && ls -la /opt

# Install service
FROM saarsec/saarctf-base
ENV SERVICENAME=saarschleife
ARG DEBIAN_FRONTEND=noninteractive
RUN apt-get update
RUN apt-get install -y dirmngr unzip openjdk-11-jdk-headless gnupg
RUN wget -qO - https://www.mongodb.org/static/pgp/server-4.2.asc | apt-key add -
RUN echo "deb http://repo.mongodb.org/apt/debian buster/mongodb-org/4.2 main" > /etc/apt/sources.list.d/mongodb-org-4.2.list
RUN apt-get update
RUN apt-get install -y mongodb-org
RUN systemctl enable mongod
RUN wget -O /tmp/gradle.zip "https://services.gradle.org/distributions/gradle-5.6.4-bin.zip"
RUN unzip -d /opt/ /tmp/gradle.zip
RUN ln -s /opt/gradle*/bin/gradle /usr/local/bin/
COPY --from=buildstage /opt /opt/
ADD ./install.sh /opt/install-service.sh
WORKDIR /opt/
RUN ./install.sh
CMD ["python3", "-u", "/usr/local/bin/systemd-replacement.py"]