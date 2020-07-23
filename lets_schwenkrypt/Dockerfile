# Build things
FROM saarsec/saarctf-base AS buildstage
ADD ./service /opt/service
ADD ./build.sh /opt/
WORKDIR /opt/
RUN ./build.sh && ls -la /opt

# Install service
FROM saarsec/saarctf-base
ENV SERVICENAME=lets_schwenkrypt
RUN apt-get update
RUN apt-get install -y python3 socat python3-pip python3-pycryptodome
COPY --from=buildstage /opt /opt/
ADD ./install.sh /opt/install-service.sh
WORKDIR /opt/
RUN ./install.sh
CMD ["python3", "-u", "/usr/local/bin/systemd-replacement.py"]
