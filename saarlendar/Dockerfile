# Build things
FROM saarsec/saarctf-base AS buildstage
RUN apt-get update
RUN apt-get install -y bash git bc python python-yaml libpcre3 libpcre3-dev zlib1g zlib1g-dev openssl libssl-dev build-essential
ADD ./service /opt/service
ADD ./build.sh /opt/
WORKDIR /opt/
RUN ./build.sh && ls -la /opt

# Install service
FROM saarsec/saarctf-base
ENV SERVICENAME=saarlendar
COPY --from=buildstage /opt /opt/
ADD ./install.sh /opt/install-service.sh
WORKDIR /opt/
RUN ./install.sh
CMD ["python3", "-u", "/usr/local/bin/systemd-replacement.py"]
