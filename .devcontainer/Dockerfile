FROM mcr.microsoft.com/devcontainers/base:ubuntu AS dev

RUN apt-get update

RUN apt-get --no-install-recommends install -y \
    libgtk2.0-0t64 libgtk-3-0t64 libgbm-dev libnotify-dev libnss3 libxss1 libasound2t64 libxtst6 xauth xvfb && \
    # clean up
    rm -rf /var/lib/apt/lists/* && \
    apt-get clean

RUN pnpm install --recursive --frozen-lockfile .