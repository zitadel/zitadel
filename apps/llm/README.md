# apps/llm

`apps/llm` contains the optional Ollama-oriented assets for the embedded risk experiment.

The API calls Ollama directly; this directory exists to hold deployment helpers such as:

- model bootstrap scripts
- sample default model choices for Compose
- documentation for the opt-in AI profile

For the first POC, the intended flow is:

1. start the `ollama` service with the Compose `ai` profile
2. pull and preload `qwen2.5:7b` through Ollama (handled automatically by `ollama-pull-model`)
3. configure ZITADEL risk mode to `observe` for the initial experiment
4. let the embedded API risk evaluator prompt bounded recent user/session signals into Ollama

The default model (`qwen2.5:7b`) is Apache 2.0 with no restrictions on automated decisions. You can switch to `enforce` mode after validating its output quality.

The simplest way to start everything locally is:

```
pnpm nx run @zitadel/compose:start-ai
```
