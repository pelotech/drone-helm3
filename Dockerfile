FROM golang:1.13

WORKDIR /src
COPY . /src/
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /src/build/drone-helm /src/cmd/drone-helm/main.go

FROM alpine/helm:3.8.1
MAINTAINER Joachim Hill-Grannec <joachim@pelo.tech>

COPY --from=0 /src/build/drone-helm /bin/drone-helm
COPY --from=0 /src/assets/kubeconfig.tpl /root/.kube/config.tpl

LABEL description="Helm 3 plugin for Drone 3"
LABEL base="alpine/helm"

ENTRYPOINT [ "/bin/drone-helm" ]
