FROM golang:alpine

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o binary

EXPOSE 8000

COPY --from=build /app/docs ./docs

ENTRYPOINT ["/app/binary"]