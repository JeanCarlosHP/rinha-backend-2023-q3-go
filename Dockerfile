FROM golang:1.21.1 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./pkg ./pkg

RUN CGO_ENABLED=0 GOOS=linux go build -o rinha-backend ./cmd/rinha/main.go

# FROM build AS test-stage
# RUN go test -v ./...

FROM scratch

WORKDIR /app

COPY --from=build ./app/rinha-backend ./

EXPOSE 8080

USER 1001

ENTRYPOINT ["./rinha-backend"]
