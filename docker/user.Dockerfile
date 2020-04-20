FROM alpine:latest

RUN apk update && apk upgrade && \
    apk --update --no-cache add tzdata && \
    mkdir /app 

WORKDIR /app 

EXPOSE 8082

COPY --from=drello-builder:latest /application/drello_user /app
COPY --from=drello-builder:latest /application/cmd/user/user_config.yaml /app

CMD /app/drello_user --config user_config.yaml