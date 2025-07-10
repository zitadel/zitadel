#!/usr/bin/env node

/**
 * Prepare script for standalone version
 * This script converts the monorepo version to a standalone version
 */

import fs from 'fs/promises';
import { execSync } from 'child_process';

const FILES_TO_REMOVE = [
    // Turbo is not needed for standalone builds since there are no workspace dependencies
    { file: 'turbo.json', backup: 'turbo.monorepo.backup.json' }
];

async function prepareStandalone() {
    console.log('🔧 Preparing standalone version...\n');

    const args = process.argv.slice(2);
    const shouldInstall = args.includes('--install');

    try {
        // Step 1: Copy package.standalone.json to package.json
        console.log('📦 Setting up package.json...');
        const packageStandaloneExists = await fs.access('package.standalone.json').then(() => true).catch(() => false);

        if (packageStandaloneExists) {
            // Backup current package.json
            await fs.copyFile('package.json', 'package.monorepo.backup.json');
            console.log('   💾 Backed up package.json → package.monorepo.backup.json');

            // Copy standalone version
            await fs.copyFile('package.standalone.json', 'package.json');
            console.log('   ✅ package.standalone.json → package.json');
        } else {
            throw new Error('package.standalone.json not found!');
        }

        // Step 2: Remove unnecessary files for standalone
        console.log('🗑️  Removing monorepo-specific files...');
        for (const item of FILES_TO_REMOVE) {
            const fileExists = await fs.access(item.file).then(() => true).catch(() => false);

            if (fileExists) {
                // Backup current file
                await fs.copyFile(item.file, item.backup);
                console.log(`   💾 Backed up ${item.file} → ${item.backup}`);

                // Remove the file
                await fs.unlink(item.file);
                console.log(`   🗑️  Removed ${item.file} (not needed in standalone)`);
            } else {
                console.log(`   ℹ️  ${item.file} not found, skipping`);
            }
        }

        // Step 3: Install dependencies if requested
        if (shouldInstall) {
            console.log('\n📥 Installing dependencies...');
            try {
                execSync('pnpm install', { stdio: 'inherit' });
                console.log('   ✅ Dependencies installed successfully');
            } catch (error) {
                console.warn('   ⚠️  pnpm install failed, you may need to run it manually');
            }
        }

        console.log('\n🎉 Standalone preparation complete!');
        console.log('   ✨ Turbo removed - using standard npm scripts');
        console.log('\n📋 Next steps:');
        if (!shouldInstall) {
            console.log('   1. Run: pnpm install');
        }
        console.log('   2. Run: pnpm run dev');
        console.log('   3. Start developing!\n');

    } catch (error) {
        console.error('\n❌ Failed to prepare standalone version:', error.message);
        console.error('Please check the error above and try again.\n');
        process.exit(1);
    }
}

prepareStandalone();
