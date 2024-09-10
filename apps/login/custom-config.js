const customConfig = {
  session: {
    lifetime_in_seconds: 7200,
  },
  selfservice: {
    change_password: {
      enabled: false,
    },
  },
};

module.exports = customConfig;
