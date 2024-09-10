import { Config } from "./default-config";

const customConfig: Partial<Config> = {
  session: {
    lifetime_in_seconds: 3600,
  },
};

export default customConfig;
