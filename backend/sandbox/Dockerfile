FROM ubuntu:22.04

# Avoid tzdata interactive prompts
ENV DEBIAN_FRONTEND=noninteractive

# Update and install common runtime dependencies
RUN apt-get update && apt-get install -y \
    curl \
    git \
    build-essential \
    python3 \
    python3-pip \
    python3-venv \
    sudo \
    jq \
    && rm -rf /var/lib/apt/lists/*

# Install Node.js (v20)
RUN curl -fsSL https://deb.nodesource.com/setup_20.x | bash - \
    && apt-get install -y nodejs \
    && npm install -g pnpm yarn

# Install Go
RUN curl -fsSL https://go.dev/dl/go1.22.2.linux-amd64.tar.gz -o go.tar.gz \
    && tar -C /usr/local -xzf go.tar.gz \
    && rm go.tar.gz
ENV PATH=$PATH:/usr/local/go/bin

# Set working directory for the agent
WORKDIR /workspace

# Keep container running
CMD ["tail", "-f", "/dev/null"]
