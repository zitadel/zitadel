import fs from "fs";
import path from "path";
import defaultConfig, { Config } from "./default-config";

const customConfigPath = path.resolve(process.cwd(), "custom-config.js");

let customConfig: Partial<Config> = {};

if (fs.existsSync(customConfigPath)) {
  import(customConfigPath)
    .then((module) => {
      customConfig = module.default;
    })
    .catch((error) => {
      console.warn("Error loading custom configuration:", error);
    });
} else {
  console.info("No custom configuration file found!");
}

const config: Config = {
  ...defaultConfig,
  ...customConfig,
};

export default config;
