FROM scratch
MAINTAINER Brian Hechinger <wonko@4amlunch.net>

ADD purge-cmd-linux-amd64 purge-cmd
VOLUME /etc/chremoas

ENTRYPOINT ["/purge-cmd", "--configuration_file", "/etc/chremoas/chremoas.yaml"]
