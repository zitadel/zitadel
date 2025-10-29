import { defineConfig } from 'vitest/config';

export default defineConfig({
  test: {
    globals: true,
    environment: 'node',
    include: ['*.test.ts', '*.spec.ts', '*.e2e.test.ts'],
    // E2E tests run sequentially to avoid port conflicts
    sequence: {
      hooks: 'list',
    },
  },
});
