FROM golang:1.17-alpine

WORKDIR /app
COPY . .

RUN go mod download

EXPOSE 8000

CMD ["go", "run", "."]