import type { Config } from '@jest/types';

export default async (): Promise<Config.InitialOptions> => {
  return {
    preset: 'ts-jest',
    transform: {
      "^.+\\.tsx?$": ['ts-jest', { tsconfig: 'tsconfig.test.json' }],
    },
    setupFilesAfterEnv: [
      '@testing-library/jest-dom/extend-expect',
    ],
    moduleNameMapper: {
      '^#/(.*)$': '<rootDir>/$1',
    },
    testEnvironment: 'jsdom',
  };
};