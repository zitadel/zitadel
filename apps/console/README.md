# ZITADEL Console (Next.js)

The next-generation ZITADEL management console, built with Next.js and React.

## Getting Started

### Prerequisites
- Node.js 20+
- pnpm
- A running ZITADEL instance with a Personal Access Token (PAT)

### Setup

1. **Copy the environment file:**
   ```bash
   cp apps/console/.env.example apps/console/.env
   ```

2. **Configure environment variables:**
   ```dotenv
   ZITADEL_INSTANCE_URL=https://your-instance.zitadel.cloud
   ZITADEL_PAT=your-personal-access-token
   ```

3. **Install dependencies** (from repo root):
   ```bash
   pnpm install
   ```

4. **Start the dev server:**
   ```bash
   pnpm nx dev console-next
   ```

   The console will be available at [http://localhost:3000](http://localhost:3000).

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Framework | Next.js 16 (App Router) |
| UI | React 19, shadcn/ui, TailwindCSS 4 |
| API | ConnectRPC → ZITADEL v2 APIs |
| Auth | Personal Access Token (PAT) |
| Icons | Lucide React |
| Monorepo | Nx + pnpm |

## Project Structure

```
apps/console/
├── app/                    # Next.js App Router pages
│   ├── users/              # User list & detail
│   ├── projects/           # Project list & detail
│   ├── applications/       # Application list & detail
│   ├── sessions/           # Session list
│   ├── organizations/      # Organization list
│   └── layout.tsx          # Root layout with providers
├── components/
│   ├── ui/                 # shadcn/ui components
│   ├── layout/             # Sidebar, account dropdown
│   └── users/              # Domain-specific components
├── lib/
│   ├── api/                # Server actions for ZITADEL APIs
│   ├── context/            # React context providers
│   ├── permissions/        # Permission gating
│   └── deployment/         # Self-hosted vs cloud features
└── hooks/                  # Custom React hooks
```

## API Architecture

All API calls go through server actions in `lib/api/`, which use ConnectRPC to communicate with ZITADEL v2 proto services:

- **Transport** (`transport.ts`): Configures ConnectRPC with the PAT bearer token.
- **Domain modules** (`users.ts`, `projects.ts`, etc.): Typed wrappers around specific v2 service RPCs.
- **Fetch helpers** (`fetch-users.ts`, etc.): JSON-safe wrappers for client component consumption.

### Adding a New API Call

1. Find the RPC in `proto/zitadel/<domain>/v2/<service>.proto`.
2. Import the request/response schemas from `@zitadel/proto`.
3. Create a `"use server"` function in the appropriate `lib/api/` file.
4. Use `create()` to build the request and `toJson()` to serialize the response.

## Available Scripts

```bash
# Development
pnpm nx dev console-next

# Production build
pnpm nx build console-next

# Lint
pnpm nx lint console-next
```

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `ZITADEL_INSTANCE_URL` | ✅ | ZITADEL instance URL (gRPC/Connect endpoint) |
| `ZITADEL_PAT` | ✅ | Personal Access Token for API authentication |
| `NEXT_PUBLIC_DEPLOYMENT_MODE` | ❌ | `self-hosted` (default) or `cloud` |
