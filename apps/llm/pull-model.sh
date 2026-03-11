#!/bin/sh
set -eu

host="${OLLAMA_HOST:-http://ollama:11434}"
model="${OLLAMA_MODEL:-qwen2.5:0.5b}"
retries="${OLLAMA_PULL_RETRIES:-60}"
sleep_seconds="${OLLAMA_PULL_SLEEP_SECONDS:-2}"

attempt=1
while [ "$attempt" -le "$retries" ]; do
  if OLLAMA_HOST="$host" ollama list >/dev/null 2>&1; then
    echo "pulling model $model ..." >&2
    env OLLAMA_HOST="$host" ollama pull "$model"

    # Preload the model into memory so the first inference request is fast.
    # keep_alive:-1 tells Ollama to keep it loaded until the server restarts.
    echo "preloading model into memory (keep_alive: -1) ..." >&2
    curl -sf -X POST "$host/api/generate" \
      -H "Content-Type: application/json" \
      -d "{\"model\":\"$model\",\"prompt\":\"\",\"keep_alive\":\"-1\",\"stream\":false}" \
      -o /dev/null
    echo "model ready." >&2
    exit 0
  fi

  echo "waiting for ollama at $host (attempt $attempt/$retries)" >&2
  attempt=$((attempt + 1))
  sleep "$sleep_seconds"
done

echo "ollama did not become ready at $host" >&2
exit 1
