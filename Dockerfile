FROM golang:1.18-alpine as builder
WORKDIR /src/
RUN apk add alpine-sdk
ADD . .
RUN cd ./cmd/drone-helm && CGO_ENABLED=0 go build -o plugin

FROM alpine/helm:3.9.3
WORKDIR /app
COPY --from=builder /src/assets/kubeconfig.tpl /root/.kube/config.tpl
COPY --from=builder /src/cmd/drone-helm/plugin ./
ENTRYPOINT [ "/app/plugin" ]
