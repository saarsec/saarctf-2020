# Build things
FROM saarsec/saarctf-base AS buildstage
ADD ./service /opt/service
ADD ./build.sh /opt/
WORKDIR /opt/
RUN ./build.sh && ls -la /opt

# Install service
FROM saarsec/saarctf-base
ENV SERVICENAME=mensaar
RUN apt-get update
RUN apt-get install -y nginx php-fpm php-pgsql python3-minimal python3-psycopg2 postgresql
COPY --from=buildstage /opt /opt/
ADD ./install.sh /opt/install-service.sh
WORKDIR /opt/
RUN ./install.sh
CMD ["python3", "-u", "/usr/local/bin/systemd-replacement.py"]
