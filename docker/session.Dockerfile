FROM alpine:latest

RUN apk update && apk upgrade && \
    apk --update --no-cache add tzdata && \
    mkdir /app 

WORKDIR /app 

EXPOSE 8081

COPY --from=drello-builder:latest /application/drello_session /app
COPY --from=drello-builder:latest /application/cmd/session/session_config.yaml /app

CMD /app/drello_session --config session_config.yaml