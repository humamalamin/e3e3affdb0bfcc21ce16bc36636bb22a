FROM golang:latest AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download && go mod verify

COPY . .

COPY --from=build /app/main /main

COPY --from=build /app/docs ./docs

ENTRYPOINT ["/main"]