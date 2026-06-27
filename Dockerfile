FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y \
    gcc g++ \
    clang \
    cmake \
    git \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY cnext-linux /usr/local/bin/cnext
RUN chmod +x /usr/local/bin/cnext

COPY test-project /app/test-project

WORKDIR /app/test-project

CMD ["cnext", "build"]
