FROM golang:1.19.5-alpine3.17 as BUILDER
RUN apk add gcc musl-dev
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY src ./src
RUN go build ./src/main.go

FROM alpine:3.17
WORKDIR /app
COPY --from=BUILDER /app/src/home.html /app/src/home.html
COPY --from=BUILDER /app/src/compose.html /app/src/compose.html
COPY --from=BUILDER /app/main .
CMD [ "./main" ]