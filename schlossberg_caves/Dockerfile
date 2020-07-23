# Build things
FROM node:lts-buster AS buildstage
ADD ./service /opt/service
ADD ./build.sh /opt/
WORKDIR /opt/
RUN ./build.sh && ls -la /opt

# Install service
FROM saarsec/saarctf-base
ENV SERVICENAME=schlossbergcaves
COPY --from=buildstage /opt /opt/
ADD ./install.sh /opt/install-service.sh
ADD ./gcc_configure.sh /opt/gcc_configure.sh
WORKDIR /opt/
RUN ./install.sh
CMD ["python3", "-u", "/usr/local/bin/systemd-replacement.py"]
