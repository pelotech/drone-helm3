FROM alpine/helm:3.8.1
MAINTAINER Erin Call <erin@liffft.com>

COPY build/drone-helm /bin/drone-helm
COPY assets/kubeconfig.tpl /root/.kube/config.tpl

LABEL description="Helm 3 plugin for Drone 3"
LABEL base="alpine/helm"

ENTRYPOINT [ "/bin/drone-helm" ]
