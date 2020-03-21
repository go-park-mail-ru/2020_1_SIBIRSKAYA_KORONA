#builder
FROM golang:alpine3.11 as builder

RUN apk update && apk upgrade && \
    apk --update add git make

WORKDIR /application

COPY . .

RUN make build-binary

#production
FROM alpine:latest as production

RUN apk update && apk upgrade && \
    apk --update --no-cache add tzdata && \
    mkdir /app 

WORKDIR /app 

EXPOSE 8080

COPY --from=builder /application/drello_binary /app

CMD /app/drello_binary