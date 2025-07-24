#!/usr/bin/env node

import fs from 'fs';
import path from 'path';
import { glob } from 'glob';

const PACKAGE_JSON_PATH = './package.json';

/**
 * Generate clean, minimal exports using wildcards for better maintainability
 * This approach is compatible with older moduleResolution: "node" strategies
 */
async function generateExports() {
    console.log('🔍 Scanning for generated proto files...');

    // Check if types directory exists (proto files have been generated)
    if (!fs.existsSync('./types')) {
        console.log('⚠️  No types directory found. Proto files may not be generated yet.');
        console.log('   Run "pnpm run generate" to generate proto files first.');
        return;
    }

    // Find all .d.ts files to validate they exist, but we don't need individual exports
    const typeFiles = await glob('types/**/*.d.ts', {
        cwd: process.cwd(),
        ignore: ['node_modules/**']
    });

    console.log(`📁 Found ${typeFiles.length} proto type files`);

    // Read current package.json
    const packageJson = JSON.parse(fs.readFileSync(PACKAGE_JSON_PATH, 'utf8'));

    // Create minimal, clean exports using wildcards
    // This is much more maintainable and works great with older module resolution
    const exports = {
        '.': {
            types: './types/zitadel/policy_pb.d.ts',
            import: './es/zitadel/policy_pb.js',
            require: './cjs/zitadel/policy_pb.cjs'
        },
        // Wildcard patterns for all proto directories
        './zitadel/*': {
            types: './types/zitadel/*.d.ts',
            import: './es/zitadel/*.js',
            require: './cjs/zitadel/*.cjs'
        },
        './zitadel/*.js': {
            types: './types/zitadel/*.d.ts',
            import: './es/zitadel/*.js',
            require: './cjs/zitadel/*.js'
        },
        './validate/*': {
            types: './types/validate/*.d.ts',
            import: './es/validate/*.js',
            require: './cjs/validate/*.cjs'
        },
        './validate/*.js': {
            types: './types/validate/*.d.ts',
            import: './es/validate/*.js',
            require: './cjs/validate/*.js'
        },
        './google/*': {
            types: './types/google/*.d.ts',
            import: './es/google/*.js',
            require: './cjs/google/*.cjs'
        },
        './google/*.js': {
            types: './types/google/*.d.ts',
            import: './es/google/*.js',
            require: './cjs/google/*.js'
        },
        './protoc-gen-openapiv2/*': {
            types: './types/protoc-gen-openapiv2/*.d.ts',
            import: './es/protoc-gen-openapiv2/*.js',
            require: './cjs/protoc-gen-openapiv2/*.cjs'
        },
        './protoc-gen-openapiv2/*.js': {
            types: './types/protoc-gen-openapiv2/*.d.ts',
            import: './es/protoc-gen-openapiv2/*.js',
            require: './cjs/protoc-gen-openapiv2/*.js'
        }
    };

    // Update package.json with clean exports
    packageJson.exports = exports;

    // Write back to package.json
    fs.writeFileSync(PACKAGE_JSON_PATH, JSON.stringify(packageJson, null, 2) + '\n');

    console.log(`✅ Generated clean exports with ${Object.keys(exports).length} patterns covering ${typeFiles.length} proto files`);
    console.log('📦 Updated package.json with minimal, maintainable exports');
    console.log(`📊 Reduced from ${typeFiles.length} explicit exports to ${Object.keys(exports).length} wildcard patterns`);
}

// Run the script
generateExports().catch(console.error);
