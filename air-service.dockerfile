FROM alpine:latest

RUN mkdir /app

COPY airApp /app

CMD ["/app/airApp"]