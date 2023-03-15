FROM golang:1.19.5-alpine3.17 as BUILDER
RUN apk add gcc musl-dev
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY cmd ./cmd
RUN go build ./cmd/main.go

FROM alpine:3.17
WORKDIR /app
COPY --from=BUILDER /app/cmd/home.html /app/cmd/home.html
COPY --from=BUILDER /app/cmd/compose.html /app/cmd/compose.html
COPY --from=BUILDER /app/main .
CMD [ "./main" ]