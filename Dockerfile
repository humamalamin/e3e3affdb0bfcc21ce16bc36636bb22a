# FROM golang:latest AS build

# ARG upx_version=4.2.4

# RUN apt-get update && apt-get install -y --no-install-recommends xz-utils && \
#   curl -Ls https://github.com/upx/upx/releases/download/v${upx_version}/upx-${upx_version}-amd64_linux.tar.xz -o - | tar xvJf - -C /tmp && \
#   cp /tmp/upx-${upx_version}-amd64_linux/upx /usr/local/bin/ && \
#   chmod +x /usr/local/bin/upx && \
#   apt-get remove -y xz-utils && \
#   rm -rf /var/lib/apt/lists/*

# # Set the working directory inside the container
# WORKDIR /app

# COPY go.mod ./
# COPY go.sum ./

# RUN go mod download && go mod verify

# # Copy the Go application files to the container
# COPY . .

# # Build the Go application
# RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o main -a -ldflags="-s -w" -installsuffix cgo

# RUN upx --ultra-brute -qq main && upx -t main

# FROM scratch

# COPY --from=build /app/main /main
# COPY --from=build /app/docs ./docs

# EXPOSE 8000

# # Set the command to run the Go application
# ENTRYPOINT ["/main"]

# Gunakan base image untuk Go
FROM golang:latest AS builder

# Set working directory di dalam container
WORKDIR /app

# Copy go.mod dan go.sum untuk mendownload dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy semua kode ke dalam container
COPY . .

# Build aplikasi Go
RUN go build -o main .

# Gunakan base image yang lebih kecil untuk runtime
FROM gcr.io/distroless/base-debian10

# Copy file biner yang sudah di build dari stage builder
COPY --from=builder /app/main /app/main
COPY --from=builder /app/docs ./docs

COPY .env /app/.env


# Set working directory
WORKDIR /app

# Expose port yang akan digunakan aplikasi (ubah sesuai aplikasi Anda)
EXPOSE 8000

# Jalankan aplikasi
CMD ["/app/main"]