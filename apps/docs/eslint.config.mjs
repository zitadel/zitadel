import { createRequire } from 'module';
const require = createRequire(import.meta.url);
const mdx = require('eslint-plugin-mdx');
const mdxParser = require('eslint-mdx');
const tsParser = require('@typescript-eslint/parser');
const tsPlugin = require('@typescript-eslint/eslint-plugin');
import globals from 'globals';

import js from '@eslint/js';
// import { FlatCompat } from '@eslint/eslintrc';
// import path from 'path';
// import { fileURLToPath } from 'url';

// const __filename = fileURLToPath(import.meta.url);
// const __dirname = path.dirname(__filename);

// const compat = new FlatCompat({
//   baseDirectory: __dirname,
// });

export default [
  // 1. Base JS Rules
  js.configs.recommended,

  // 2. Import Next.js Core Web Vitals (via compat because it doesn't support Flat Config natively yet)
  // ...compat.extends('next/core-web-vitals'),
  // ...compat.extends('next', 'next/typescript'), // Causes circular JSON error


  // ...
  // 3. MDX Configuration
  {
    files: ['*.mdx'],
    ...mdx.flat,
    processor: mdx.createRemarkProcessor({
      lintCodeBlocks: false,
    }),
    languageOptions: {
      globals: {
        ...globals.browser,
        ...globals.node,
        props: 'readonly',
      },
      parser: mdxParser,
      ecmaVersion: 'latest',
      sourceType: 'module',
    },
    rules: {
      'react/no-unescaped-entities': 'off', 
      'no-unused-vars': 'off',
      'no-irregular-whitespace': 'off',
      'no-useless-escape': 'off',
    },
  },
  // ...
  
  // 4. Code Blocks inside MDX (virtual files)
  {
    files: ['*.mdx/*.{js,jsx,ts,tsx}'],
    languageOptions: {
      globals: {
        ...globals.browser,
        ...globals.node,
        props: 'readonly', // common in MDX
      }
    },
    rules: {
      'no-unused-vars': 'off',
      'import/no-anonymous-default-export': 'off',
    },
  },
  
  // 5. Global Ignores
  {
    ignores: ["node_modules/", ".next/", "out/", "build/", ".source/", "next-env.d.ts"],
  },
  
  // 6. Global options for JS files
  {
    files: ['**/*.{js,mjs,cjs}'],
    languageOptions: {
      globals: {
        ...globals.browser,
        ...globals.node,
      },
    }
  },
  // 7. Manual TypeScript Configuration
  {
    files: ['**/*.{ts,tsx}'],
    plugins: {
      '@typescript-eslint': tsPlugin,
    },
    languageOptions: {
      parser: tsParser,
      parserOptions: {
        ecmaVersion: 'latest',
        sourceType: 'module',
        ecmaFeatures: { jsx: true },
      },
      globals: {
        ...globals.browser,
        ...globals.node,
      },
    },
    rules: {
      ...tsPlugin.configs.recommended.rules,
      'no-unused-vars': 'off',
      '@typescript-eslint/no-unused-vars': ['warn', { argsIgnorePattern: '^_' }],
      '@typescript-eslint/no-explicit-any': 'off', // Disable for legacy code
      '@typescript-eslint/ban-ts-comment': 'off',
    },
  },
  // 8. Scripts Config
  {
    files: ['scripts/**/*.{js,mjs,ts}'],
    rules: {
      'no-unused-vars': 'off',
      '@typescript-eslint/no-unused-vars': 'off',
      'no-undef': 'off',
      'no-useless-escape': 'off',
      'no-empty': 'off',
      'no-async-promise-executor': 'off',
    }
  }
];
