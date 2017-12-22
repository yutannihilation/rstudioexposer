FROM golang:alpine as builder
COPY login.go main.go /tmp/
WORKDIR /tmp/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o rstudioexposer .

FROM rocker/tidyverse:latest
COPY --from=builder /tmp/rstudioexposer /usr/local/bin/rstudioexposer
RUN mkdir -p /etc/services.d/rstudioexposer
COPY files/run  /etc/services.d/rstudioexposer/run

EXPOSE 80 8787