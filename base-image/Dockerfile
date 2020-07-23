FROM debian:buster

RUN apt-get update && \
	apt-get upgrade -y && \
	apt-get install -y --no-install-recommends wget ca-certificates sudo nano htop && \
	apt-get install -y python3 python3-pip python3-setuptools python3-wheel python3-virtualenv && \
	apt-get clean && \
	true

ADD usr /usr/
RUN chmod -R +x /usr/local/bin/
RUN bash -c 'rm -f /lib/systemd/system/multi-user.target.wants/* \
  /etc/systemd/system/*.wants/* \
  /lib/systemd/system/local-fs.target.wants/* \
  /lib/systemd/system/sockets.target.wants/*udev* \
  /lib/systemd/system/sockets.target.wants/*initctl* \
  /lib/systemd/system/sysinit.target.wants/systemd-tmpfiles-setup* \
  /lib/systemd/system/systemd-update-utmp*'
ADD opt /opt/
