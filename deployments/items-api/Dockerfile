############################
# STEP 1 build executable binary
############################
FROM golang:1.18-buster as builder

ENV GO111MODULE=on
WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/server /app/cmd/items-api

############################
# STEP 2 build a small image
############################
FROM alpine:3.11.6

WORKDIR /app

# Copy only required files to final docker image
COPY --from=builder /app/Makefile  /app/Makefile
COPY --from=builder /app/bin/server /app/bin/server
COPY --from=builder /app/.env /app/.env
COPY --from=builder /app/migration/* /app/migration/
COPY --from=builder /app/app/items-api/api/items.yml /app/static/swagger.yml

RUN chmod +x /app/bin/server
EXPOSE 8080
ENTRYPOINT ["/app/bin/server"]
