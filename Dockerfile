# syntax=docker/dockerfile:1

FROM ghcr.io/merliot/device/device-base:latest

WORKDIR /app
RUN git clone https://github.com/merliot/device.git
RUN go work use device

WORKDIR /app/garage

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go work use .
RUN go build -tags prime -o /garage ./cmd/
RUN go run ../device/cmd/uf2-builder -target nano-rp2040 -model garage
RUN go run ../device/cmd/uf2-builder -target wioterminal -model garage

EXPOSE 8000

ENV PORT_PRIME=8000
CMD ["/garage"]
