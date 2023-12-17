
FROM golang:latest


WORKDIR /app

COPY . .


RUN go build -o app/file

EXPOSE 8000

CMD ["./app/file"]
