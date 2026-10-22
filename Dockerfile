FROM golang:1.27-trixie

# Basic tools
RUN apt-get update && apt-get install -y --no-install-recommends \
    git \
    curl \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

RUN curl -fsSL https://antigravity.google/cli/install.sh | bash

# Go development environment
ENV GOPATH=/go
ENV PATH="/go/bin:/usr/local/go/bin:${PATH}"

WORKDIR /workspace
