FROM golang:latest AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download && go mod verify

COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o main


COPY --from=build /app/main /main
COPY --from=build /app/docs ./docs

# Set the command to run the Go application
ENTRYPOINT ["/main"]