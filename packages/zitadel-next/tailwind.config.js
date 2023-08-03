const sharedConfig = require("zitadel-tailwind-config/tailwind.config.js");

/** @type {import('tailwindcss').Config} */
module.exports = {
  presets: [sharedConfig],
  prefix: "ztdl-next-",
  darkMode: "class",
  content: [`src/**/*.{js,ts,jsx,tsx}`],
};
