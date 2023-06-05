import type { Config } from '@jest/types';
import { pathsToModuleNameMapper } from 'ts-jest'
import { compilerOptions } from './tsconfig.json';

export default async (): Promise<Config.InitialOptions> => {
  return {
    preset: 'ts-jest',
    transform: {
      "^.+\\.tsx?$": ['ts-jest', {tsconfig:'./__test__/tsconfig.json'}],
    },
    setupFilesAfterEnv: [
      '@testing-library/jest-dom/extend-expect',
    ],
    moduleNameMapper: pathsToModuleNameMapper(compilerOptions.paths, {
      prefix:'<rootDir>/'
    }),
    testEnvironment: 'jsdom',
    testRegex: '/__test__/.*\\.test\\.tsx?$',
    modulePathIgnorePatterns: ['cypress'],
    reporters: [[ 'jest-silent-reporter', { useDots: true, showWarnings: true, showPaths: true } ], 'summary']
  };
};