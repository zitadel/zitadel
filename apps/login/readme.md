# ZITADEL Login UI

This is going to be our next UI for the hosted login. It's based on Next.js 13 and its introduced `app/` directory.

## Custom Configuration

You can overwrite the default configuration by creating a `custom-config.ts` file in the root of the project. The `custom-config.ts` file should contain the settings you want to overwrite.

### Example `custom-config.ts`

```js
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
```
