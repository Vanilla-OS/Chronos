FROM ghcr.io/vanilla-os/pico:main
ENV DEBIAN_FRONTEND=noninteractive
RUN apt update && apt install git golang -y

WORKDIR /app
COPY . .

RUN go build -o chronos .

ENTRYPOINT ["./chronos"]
