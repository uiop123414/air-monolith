FROM alpine:latest

RUN mkdir /app

COPY airApp /app

COPY /internal/schemas /app/schemas

CMD ["/app/airApp"]

