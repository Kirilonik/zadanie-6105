FROM golang:1.23.1

WORKDIR /app
COPY backend/. ./
RUN go mod download -json
COPY . ./

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]