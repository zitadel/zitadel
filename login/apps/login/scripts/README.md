# Standalone Build Scripts

This directory contains the simplified scripts needed for managing the ZITADEL Login UI standalone conversion.

## 📁 File Overview

### Required Files

- `package.json` - Main package file (monorepo mode with `workspace:*`)
- `package.standalone.json` - Pre-configured standalone version (uses `latest`)
- `scripts/prepare-standalone.js` - Conversion script

### Generated/Temporary Files

- `package.monorepo.backup.json` - Backup when switching to standalone mode

## 🛠️ Scripts Overview

### `prepare-standalone.js`

**The main script for converting monorepo to standalone mode.**

```bash
node scripts/prepare-standalone.js [--install]
```

- `--install` - Automatically install dependencies after preparation

**What it does:**

- Copies `package.standalone.json` → `package.json`
- Confirms all configurations are unified and ready for standalone use
- Optionally runs `pnpm install`

## 🚀 **Simplified Approach**

This setup now uses a **much simpler approach**:

1. **Unified Configuration**: All ESLint, Prettier, and Tailwind configs work for both monorepo and standalone modes
2. **Static Configuration**: `package.standalone.json` is pre-configured with `latest` versions
3. **No Duplicate Configs**: No separate `*.standalone.*` config files needed
4. **Faster Setup**: Conversion is instant with just file copying

## 📋 **Usage for Customers**

```bash
# 1. Clone the repository
git clone <standalone-repo>

# 2. Prepare standalone mode
node scripts/prepare-standalone.js --install

# 3. Start developing
pnpm run dev
```

## 🔧 **Usage for Maintainers**

```bash
# Update to latest packages manually in package.standalone.json
npm view @zitadel/client version
npm view @zitadel/proto version

# Then update package.standalone.json with latest versions
```

## ✨ **Key Benefits**

- **Single Config Files**: ESLint, Prettier, and Tailwind configs work for both modes
- **No Duplication**: No need for separate `*.standalone.*` configuration files
- **Faster Conversion**: Only 2 files need to be copied (package.json and tsconfig.json)
- **Simpler Maintenance**: All configuration logic is in one place

# Switch back to monorepo mode

node scripts/config-manager.js switch monorepo

````

### `validate-standalone.js`

**Validates that standalone setup is working correctly.**

```bash
node scripts/validate-standalone.js
````

**Checks:**

- Required files exist
- Package.json has required scripts and dependencies
- Dependencies can be resolved
- TypeScript compilation works

## Workflow Examples

### For Developers (Monorepo)

1. **Working in monorepo:**
