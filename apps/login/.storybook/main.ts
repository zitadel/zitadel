import type { StorybookConfig } from "@storybook/nextjs";
const config: StorybookConfig = {
  stories: [ "../**!(node_modules)/*.stories.@(js|jsx|ts|tsx|md|mdx)"],
  addons: [
    "@storybook/addon-links",
    "@storybook/addon-essentials",
    "@storybook/addon-interactions",
    'storybook-css-modules-preset',
  ],
  framework: {
    name: "@storybook/nextjs",
    options: {},
  },
  docs: {
    autodocs: "tag",
  },
  core: {
    disableTelemetry: true,
  },

};
export default config;
