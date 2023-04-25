const sharedConfig = require("zitadel-tailwind-config/tailwind.config.js");

/** @type {import('tailwindcss').Config} */
module.exports = {
  presets: [sharedConfig],
  prefix: "ztdl-",
  darkMode: "class",
  content: [`src/**/*.{js,ts,jsx,tsx}`],
};
