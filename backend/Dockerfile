FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod ./
COPY . .
RUN go test ./... && go build -o /calculator ./cmd/server

FROM alpine:3.20
WORKDIR /app
COPY --from=build /calculator /app/calculator
EXPOSE 8080
CMD ["/app/calculator"]
