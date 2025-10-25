/// <reference types='vitest' />
import {defineConfig} from 'vite';
import angular from '@analogjs/vite-plugin-angular';
import {nxViteTsPaths} from '@nx/vite/plugins/nx-tsconfig-paths.plugin';
import {nxCopyAssetsPlugin} from '@nx/vite/plugins/nx-copy-assets.plugin';
import {resolve} from 'node:path'

export default defineConfig(() => ({
  root: __dirname,
  cacheDir: '../node_modules/.vite/console',
  plugins: [angular(), nxViteTsPaths(), nxCopyAssetsPlugin(['*.md'])],
  // Uncomment this if you are using workers.
  // worker: {
  //  plugins: [ nxViteTsPaths() ],
  // },
  test: {
    name: 'console',
    watch: false,
    globals: true,
    environment: 'jsdom',
    include: ['{src,tests}/**/*.{test,spec}.{js,mjs,cjs,ts,mts,cts,jsx,tsx}'],
    setupFiles: ['src/test-setup.ts'],
    reporters: ['default'],
    coverage: {
      reportsDirectory: '../coverage/console',
      provider: 'v8' as const,
    },
    alias: [
      {find: "@", replacement: resolve(__dirname, "./src/app")},
      {find: "src/app", replacement: resolve(__dirname, "./src/app")},
    ]
  },
}));
