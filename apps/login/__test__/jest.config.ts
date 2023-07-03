import type { Config } from "@jest/types";
import { pathsToModuleNameMapper } from "ts-jest";
import { compilerOptions } from "../tsconfig.json";

// We make these type imports explicit, so IDEs with their own typescript engine understand them, too.
import type {} from "@testing-library/jest-dom";

export default async (): Promise<Config.InitialOptions> => {
  return {
    preset: "ts-jest",
    transform: {
      "^.+\\.tsx?$": ["ts-jest", { tsconfig: "<rootDir>/tsconfig.json" }],
    },
    setupFilesAfterEnv: ["@testing-library/jest-dom/extend-expect"],
    moduleNameMapper: pathsToModuleNameMapper(compilerOptions.paths, {
      prefix: "<rootDir>/../",
    }),
    testEnvironment: "jsdom",
    testRegex: "/__test__/.*\\.test\\.tsx?$",
    modulePathIgnorePatterns: ["cypress"],
  };
};
