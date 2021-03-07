FROM golang:alpine AS builder
WORKDIR /local/go/src/
RUN apk --no-cache add ca-certificates git
COPY ./code/ .
RUN pwd && ls
RUN go get -v && CGO_ENABLED=0 go build -o /bin/octoprint_exporter octoprint_exporter.go
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /bin/octoprint_exporter /bin/octoprint_exporter
EXPOSE 8081
ENV OCTOPRINT_API_KEY=""
ENV OCTOPRINT_HOST=""
CMD ["/bin/octoprint_exporter","-host","OCTOPRINT_HOST","-apikey","OCTOPRINT_API_KEY"]
