FROM alpine:latest

RUN apk update && apk upgrade && \
    apk --update --no-cache add tzdata && \
    mkdir /app 

WORKDIR /app 

EXPOSE 8080

COPY --from=drello-builder:latest /application/drello_api /app
COPY --from=drello-builder:latest /application/cmd/api/api_config.yaml /app

CMD /app/drello_api --config api_config.yaml