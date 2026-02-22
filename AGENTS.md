# ZITADEL Monorepo Guide for AI Agents

## Mission & Context
ZITADEL is an open-source Identity Management System (IAM) written in Go and Angular/React. It provides secure login, multi-tenancy, and audit trails.

## Read Order
1. Read this file first.
2. Read the nearest scoped `AGENTS.md` for the area you edit.
3. If multiple scopes apply, use the most specific path.

## Repository Structure Map
- **`apps/`**: Consumer-facing web applications.
  - **`login`**: Next.js authentication UI. See `apps/login/AGENTS.md`.
  - **`docs`**: Fumadocs documentation app. See `apps/docs/AGENTS.md`.
  - **`api`**: Backend Nx app target. See `apps/api/AGENTS.md`.
- **`console/`**: Angular Management Console. See `console/AGENTS.md`.
- **`internal/`**: Backend domain and service logic. See `internal/AGENTS.md`.
- **`proto/`**: API definitions. See `proto/AGENTS.md`.
- **`packages/`**: Shared TypeScript packages. See `packages/AGENTS.md`.
- **`tests/functional-ui/`**: Cypress functional UI tests. See `tests/functional-ui/AGENTS.md`.

## Technology Stack & Conventions
- **Orchestration**: Nx is used for build and task orchestration.
- **Package Manager**: pnpm.
- **Backend**:
  - **Go Version Source of Truth**: Inspect `go.mod` before Go work (`go` and optional `toolchain` directives).
  - **Communication**: For V2 APIs, connectRPC is the primary transport. gRPC and HTTP/JSON endpoints are also supported.
  - **Pattern**: The backend is transitioning to a relational design. Events are still persisted in a separate table for history/audit, but events are not the system of record.
- **Frontend**:
  - **Console**: Angular + RxJS.
  - **Login/Docs**: Next.js + React.

## ZITADEL Domain & Multi-Tenancy Logic

### 1. Hierarchy & Ownership

ZITADEL follows a strict hierarchical containment model. When generating code, logic, or translations, adhere to this structure:

- **System (Installation):** The entire ZITADEL deployment. Global settings are applied through runtime configuration files or environment variables. See `cmd/defaults.yaml`.
- **Instance (The "Identity System"):**
  - **Definition:** A logical partition/virtual tenant. It is a "System inside a System."
  - **Isolation:** Data and settings are strictly isolated between instances.
  - **Translation Rule:** NEVER translate as "Example" or "Case." Use technical terms like "Tenant," "Environment," or the local equivalent of "Logical System Entity."
- **Organization:** A group within an Instance. It owns Users, Projects, and Roles.
- **Project:** A collection of Applications and Auth Policies within an Org.

### 2. Permission Scoping (The "Administrative" Context)

- **System User:** Manages the entire Installation and creates Instances over the system API.
- **Instance Admin:** Manages Instance-wide policies (Password complexity, Identity Providers, Organizations).
- **Organization Admin:** Manages users and access within a specific Organization.

### 3. Language & Tone Guidelines

- **Avoid Ambiguity:** When referring to an 'Instance', the context is always infrastructure/tenancy.
- **Technical Precision:** In UI strings, prefer clarity over brevity if "Instance" is likely to be misinterpreted in the target language.

### 4. Technical Glossary & Localization Mapping

| Language             | Technical Term (SaaS/Cloud) | Why this term?                                            | "Avoid this (The ""Example"" Trap)" |
| -------------------- | --------------------------- | --------------------------------------------------------- | ----------------------------------- |
| Chinese (Simplified) | 实例 (Shílì)                | Standard for a cloud resource/entity.                     | 例子 (Lìzi)                         |
| Japanese             | インスタンス (Insutansu)    | Katakana transliteration; industry standard.              | 例 (Rei)                            |
| Korean               | 인스턴스 (Inseuteonseu)     | Hangul transliteration; industry standard.                | 예 (Ye)                             |
| German               | Instanz                     | Matches English but implies a technical occurrence.       | Beispiel                            |
| French               | Instance                    | "Standard but often requires ""de ZITADEL"" for clarity." | Exemple                             |
| Spanish              | Instancia                   | Technical entity in software architecture.                | Ejemplo                             |
| Portuguese           | Instância                   | Standard technical terminology.                           | Exemplo                             |
| Russian              | Инстанс (Instans)           | Modern SaaS jargon (transliterated).                      | Пример (Primer)                     |

#### Translation Guardrails
If a translation is requested for a language not listed above, follow these priority rules for the word 'Instance':

1. **Priority 1 (Transliteration):** Use the phonetic transliteration into the local script (common in Japanese/Korean/Russian).
2. **Priority 2 (System Entity):** Use a term that implies a "running process" or "logical environment."
3. **Priority 3 (Tenant):** If 'Instance' is ambiguous, use the local word for 'Tenant' (e.g., 租户 in Chinese).
4. **Strict Ban:** NEVER use words that mean "an illustration", "a case", "a sample", or "an example."

### 5. Deployment Targets

ZITADEL supports multiple deployment methods. Each has its own directory and conventions:

| Target | Location | Status | Notes |
|--------|----------|--------|-------|
| Docker Compose | `deploy/compose/` | Supported | Single-node, graduated from quickstart to semi-production. See `deploy/compose/agents.md` for directory-specific rules. |
| Kubernetes (Helm) | External ([zitadel-charts](https://github.com/zitadel/zitadel-charts)) | Supported | Official Helm chart for production workloads. Docs at `apps/docs/content/self-hosting/deploy/kubernetes/`. |
| apt/rpm packages | Planned | Not yet available | Future packaging target. |

When generating deployment-related content:
- Docker Compose is the recommended path for **getting started** and **homelab/single-node** deployments
- Kubernetes is the recommended path for **production** workloads
- Always reference the correct deployment method for the user's context
- The same `ZITADEL_*` environment variable model applies across all deployment methods

## Command Rules
Run commands from the repository root.

- Use verified Nx targets only.
- If target availability is unclear, run `pnpm nx show project <project>`.
- Do not assume all projects have `test`, `lint`, `build`, or `generate` targets.
- Known exception: `@zitadel/console` has no configured `test` target.

## Verified Common Targets
- `@zitadel/api`: `prod`, `build`, `generate`, `generate-install`, `lint`, `test`, `test-unit`, `test-integration`
- `@zitadel/login`: `dev`, `build`, `lint`, `test`, `test-unit`, `test-integration`
- `@zitadel/docs`: `dev`, `build`, `generate`, `install-proto-plugins`, `check-links`, `check-types`, `test`, `lint`
- `@zitadel/console`: `dev`, `build`, `generate`, `install-proto-plugins`, `lint`
- `@zitadel/compose`: `test-config`, `test-run`, `test-e2e`, `test`, `stop`, `test-login-acceptance`

## Proto Plugin Binaries
All proto plugins are installed to `.artifacts/bin/<GOOS>/<GOARCH>/` and Nx-cached. `generate` targets wire up the correct install dependency and prepend `.artifacts/bin/` to `$PATH` — no manual install step is needed.

## PR Title Convention

PR titles are validated by the Semantic PR app. Format:

`<type>(<scope>): <short summary>`

**Types**: must come from the list in [`.github/semantic.yml`](.github/semantic.yml) under `types:` — e.g. `feat`, `fix`, `docs`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`, `revert`.

**Scopes**: optional, but if used must come from the list in [`.github/semantic.yml`](.github/semantic.yml) under `scopes:`. When in doubt, omit the scope — do not invent values not on that list.

## Documentation
- **Human Guide**: See `CONTRIBUTING.md` for setup and contribution details.
- **API Design**: See `API_DESIGN.md` for API specific guidelines.
