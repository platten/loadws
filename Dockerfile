# build stage
FROM golang:alpine AS build-env

ADD https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep
RUN apk update; apk add git

WORKDIR $GOPATH/src/github.com/platten/loadWS
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only

COPY loadws.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /loadws .

# final stage
FROM alpine
EXPOSE 8080

WORKDIR /app
COPY --from=build-env /loadws /app/
ENTRYPOINT ./loadws