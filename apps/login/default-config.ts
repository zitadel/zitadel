export interface Config {
  session: {
    lifetime_in_seconds: number;
  };
  selfservice: {
    change_password: {
      enabled: boolean;
    };
  };
}

const defaultConfig: Config = {
  session: {
    lifetime_in_seconds: 3600,
  },
  selfservice: {
    change_password: {
      enabled: false,
    },
  },
};

export default defaultConfig;
