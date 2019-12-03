FROM alpine/helm
MAINTAINER Erin Call <erin@liffft.com>

COPY build/drone-helm /bin/drone-helm

LABEL description="Helm 3 plugin for Drone 3"
LABEL base="alpine/helm"

ENTRYPOINT [ "/bin/drone-helm" ]
