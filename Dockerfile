# build stage
FROM golang:alpine AS build-env
ADD . /src
RUN apk update && apk upgrade && \
        apk add --no-cache bash git openssh
RUN go get github.com/phanoix/graphql_usr_srv
RUN cd /src && go build -o usr_srv

#srv img
FROM alpine
WORKDIR /app
COPY --from=build-env /src/usr_srv /app/
EXPOSE 8000
ENTRYPOINT  ./usr_srv