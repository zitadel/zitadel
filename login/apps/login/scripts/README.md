# Standalone Build Scripts

This directory contains the minimal scripts needed for managing the ZITADEL Login UI standalone conversion.

## 📁 File Overview

### Required Files

- `package.json` - Main package file (monorepo mode with `workspace:*`)
- `package.standalone.json` - Pre-configured standalone version (uses `latest`)
- `*.standalone.*` - Configuration templates for standalone mode
- `scripts/` - Conversion scripts

### Generated/Temporary Files

- `*.generated.*` - Temporary files (ignored by git)
- `package.monorepo.backup.json` - Backup when switching modes

## 🛠️ Scripts Overview

### `prepare-standalone.js`

**The main script for converting monorepo to standalone mode.**

```bash
node scripts/prepare-standalone.js [--install]
```

- `--install` - Automatically install dependencies after preparation

**What it does:**

- Copies `package.standalone.json` → `package.json`
- Copies all `*.standalone.*` config files to their active versions
- Optionally runs `npm install`

## 🚀 **Simplified Approach**

This setup now uses a **much simpler approach**:

1. **Static Configuration**: `package.standalone.json` is pre-configured with `latest` versions
2. **No Version Fetching**: No network calls or complex version management
3. **Manual Maintenance**: Package versions are updated manually when needed
4. **Faster Setup**: Conversion is instant with just file copying

## 📋 **Usage for Customers**

```bash
# 1. Clone the repository
git clone <standalone-repo>

# 2. Prepare standalone mode
./prepare-standalone.sh --install

# 3. Start developing
npm run dev
```

## 🔧 **Usage for Maintainers**

```bash
# Update to latest packages manually in package.standalone.json
npm view @zitadel/client version
npm view @zitadel/proto version

# Then update package.standalone.json with latest versions
```

**What it does:**

- Converts `workspace:*` dependencies to published package versions
- Removes monorepo-specific devDependencies
- Copies standalone configuration files
- Updates package.json scripts for standalone use

### `build-standalone.js`

**Generates package.standalone.json with latest published package versions.**

```bash
node scripts/build-standalone.js
```

**What it does:**

- Fetches latest versions of `@zitadel/client` and `@zitadel/proto`
- Creates a standalone-ready package.json
- Safe to run in monorepo (doesn't modify current package.json)

### `config-manager.js`

**Utility for switching between monorepo and standalone configurations.**

```bash
# Show current mode
node scripts/config-manager.js status

# Switch to standalone mode
node scripts/config-manager.js switch standalone

# Switch back to monorepo mode
node scripts/config-manager.js switch monorepo
```

### `validate-standalone.js`

**Validates that standalone setup is working correctly.**

```bash
node scripts/validate-standalone.js
```

**Checks:**

- Required files exist
- Package.json has required scripts and dependencies
- Dependencies can be resolved
- TypeScript compilation works

## Workflow Examples

### For Developers (Monorepo)

1. **Working in monorepo:**

   ```bash
   # Normal development - uses workspace:* dependencies
   pnpm dev
   ```

2. **Testing standalone version:**

   ```bash
   # Generate standalone package.json (safe, doesn't modify current setup)
   node scripts/build-standalone.js

   # Switch to standalone mode for testing
   node scripts/config-manager.js switch standalone
   npm install
   npm run dev

   # Switch back to monorepo mode
   node scripts/config-manager.js switch monorepo
   ```

### For Customers (Standalone)

1. **Initial setup:**

   ```bash
   git clone <standalone-repo>
   cd login
   ./prepare-standalone.sh --install
   ```

2. **Development:**

   ```bash
   npm run dev
   ```

3. **Production:**
   ```bash
   npm run build:standalone
   npm run start
   ```

### For CI/CD (Package Publishing)

1. **After publishing new @zitadel packages:**

   ```bash
   # Update standalone version with latest packages
   node scripts/build-standalone.js

   # Commit and push to standalone branch/repo
   git add package.standalone.json
   git commit -m "Update standalone package versions"
   git push origin standalone
   ```

## Configuration Files

The scripts manage these configuration files:

- `package.standalone.json` - Standalone version of package.json
- `tsconfig.standalone.json` - TypeScript config for standalone
- `.eslintrc.standalone.cjs` - ESLint config for standalone
- `prettier.config.standalone.mjs` - Prettier config for standalone
- `tailwind.config.standalone.mjs` - Tailwind config for standalone

## File Structure

```
scripts/
├── prepare-standalone.js     # Main preparation script
├── build-standalone.js       # Package.json generator
├── config-manager.js         # Configuration switcher
└── README.md                # This file

# Configuration files
├── package.json              # Current active package.json
├── package.standalone.json   # Generated standalone version
├── tsconfig.json            # Current active TypeScript config
├── tsconfig.standalone.json  # Standalone TypeScript config
├── .eslintrc.cjs            # Current active ESLint config
├── .eslintrc.standalone.cjs  # Standalone ESLint config
└── ...
```

## Troubleshooting

### "Could not fetch latest version" warnings

This is normal when packages haven't been published yet or npm registry is slow. The scripts will fall back to existing versions.

### Configuration file not found

Some configuration files are optional. The scripts will warn but continue without them.

### Permission denied on scripts

Make sure scripts are executable:

```bash
chmod +x scripts/*.js
chmod +x prepare-standalone.sh
```

### Package version conflicts

If you encounter version conflicts, try:

```bash
# Clean install
rm -rf node_modules package-lock.json
npm install
```

## Integration with Monorepo

These scripts are designed to work seamlessly with the existing monorepo structure:

- **Development**: Use normal `pnpm` commands in monorepo
- **Testing**: Use `config-manager.js` to test standalone mode
- **Publishing**: Use `build-standalone.js` to generate updated standalone configs
- **Customer Distribution**: Use `prepare-standalone.sh` for easy setup
