# Build things
FROM mcr.microsoft.com/dotnet/core/sdk:3.1 AS buildstage
ADD ./service /opt/service
ADD ./build.sh /opt/
WORKDIR /opt/
RUN ./build.sh && ls -la /opt

# Install service
FROM saarsec/saarctf-base
ENV SERVICENAME=saarxiv
RUN bash -c '(wget -qO- https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor > /etc/apt/trusted.gpg.d/microsoft.asc.gpg) &&\
 wget -q https://packages.microsoft.com/config/debian/10/prod.list -O/etc/apt/sources.list.d/microsoft-prod.list &&\
 chown root:root /etc/apt/trusted.gpg.d/microsoft.asc.gpg &&\
 chown root:root /etc/apt/sources.list.d/microsoft-prod.list &&\
 apt-get update &&\
 apt-get install -y apt-transport-https aspnetcore-runtime-3.1 &&\
 apt-get install -y --no-install-recommends texlive-base texlive-latex-base'
COPY --from=buildstage /opt /opt/
ADD ./install.sh /opt/install-service.sh
WORKDIR /opt/
RUN ./install.sh
CMD ["python3", "-u", "/usr/local/bin/systemd-replacement.py"]
