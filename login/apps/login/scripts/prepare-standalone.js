#!/usr/bin/env node

/**
 * Prepare script for standalone version
 * This script converts the monorepo version to a standalone version
 */

import fs from 'fs/promises';
import { execSync } from 'child_process';

const CONFIG_FILES = [
    {
        source: 'tsconfig.standalone.json',
        target: 'tsconfig.json',
        required: true
    }
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

        // Step 2: Copy TypeScript configuration
        console.log('\n⚙️  Setting up TypeScript configuration...');
        for (const config of CONFIG_FILES) {
            try {
                const sourceExists = await fs.access(config.source).then(() => true).catch(() => false);
                if (sourceExists) {
                    await fs.copyFile(config.source, config.target);
                    console.log(`   ✅ ${config.source} → ${config.target}`);
                } else if (config.required) {
                    throw new Error(`Required file ${config.source} not found!`);
                } else {
                    console.log(`   ⚠️  ${config.source} not found, skipping`);
                }
            } catch (error) {
                if (config.required) {
                    throw error;
                }
                console.warn(`   ❌ Failed to copy ${config.source}: ${error.message}`);
            }
        }

        // Step 3: Install dependencies if requested
        if (shouldInstall) {
            console.log('\n📥 Installing dependencies...');
            try {
                execSync('npm install', { stdio: 'inherit' });
                console.log('   ✅ Dependencies installed successfully');
            } catch (error) {
                console.warn('   ⚠️  npm install failed, you may need to run it manually');
            }
        }

        console.log('\n🎉 Standalone preparation complete!');
        console.log('\n📋 Next steps:');
        if (!shouldInstall) {
            console.log('   1. Run: npm install');
        }
        console.log('   2. Run: npm run dev');
        console.log('   3. Start developing!\n');

        console.log('ℹ️  Note: ESLint, Prettier, and Tailwind configs are now unified');
        console.log('   - No separate standalone config files needed!');

    } catch (error) {
        console.error('\n❌ Failed to prepare standalone version:', error.message);
        console.error('Please check the error above and try again.\n');
        process.exit(1);
    }
}

prepareStandalone();
