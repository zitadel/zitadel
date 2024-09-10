import customConfig from "./custom-config";
import defaultConfig, { Config } from "./default-config";

const config: Config = {
  ...defaultConfig,
  ...customConfig,
};

export default config;
