# build stage
FROM golang:1.15 as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

# runtime stage
FROM scratch

WORKDIR /app

COPY ./data ./data

COPY --from=builder /app/app ./mrt

ENTRYPOINT [ "/app/mrt" ]