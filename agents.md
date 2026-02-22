## ZITADEL Domain & Multi-Tenancy Logic

### 1. Hierarchy & Ownership

ZITADEL follows a strict hierarchical containment model. When generating code, logic, or translations, adhere to this structure:

- **System (Installation):** The entire ZITADEL deployment. Global settings are applied through runtime configuration files or environment variables. See `cmd/defaults.yaml`.
- **Instance (The "Identity System"):** 
  - **Definition:** A logical partition/virtual tenant. It is a "System inside a System."
  - **Isolation:** Data and configurations are strictly isolated between instances.
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
