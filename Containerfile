FROM ghcr.io/vanilla-os/pico:main
ENV DEBIAN_FRONTEND=noninteractive
RUN apt update && apt install git golang -y

WORKDIR /app

RUN wget https://github.com/Vanilla-OS/Chronos/archive/refs/heads/main.zip && \
    unzip main.zip && \
    mv Chronos-main/* . && \
    rm -r Chronos-main main.zip

RUN go build -o chronos .

ENTRYPOINT ["./chronos"]
