const sharedConfig = require("zitadel-tailwind-config/tailwind.config.js");

let colors = {
  background: { light: { contrast: {} }, dark: { contrast: {} } },
  primary: { light: { contrast: {} }, dark: { contrast: {} } },
  warn: { light: { contrast: {} }, dark: { contrast: {} } },
  text: { light: { contrast: {} }, dark: { contrast: {} } },
  link: { light: { contrast: {} }, dark: { contrast: {} } },
};
const shades = [
  "50",
  "100",
  "200",
  "300",
  "400",
  "500",
  "600",
  "700",
  "800",
  "900",
];
const themes = ["light", "dark"];
const types = ["background", "primary", "warn", "text", "link"];
types.forEach((type) => {
  themes.forEach((theme) => {
    shades.forEach((shade) => {
      colors[type][theme][shade] = `var(--theme-${theme}-${type}-${shade})`;
      colors[type][theme][
        `contrast-${shade}`
      ] = `var(--theme-${theme}-${type}-contrast-${shade})`;
    });
  });
});

/** @type {import('tailwindcss').Config} */
module.exports = {
  presets: [sharedConfig],
  darkMode: "class",
  content: [
    "./app/**/*.{js,ts,jsx,tsx}",
    "./page/**/*.{js,ts,jsx,tsx}",
    "./ui/**/*.{js,ts,jsx,tsx}",
  ],
  future: {
    hoverOnlyWhenSupported: true,
  },
  theme: {
    extend: {
      colors,
    },
  },
  plugins: [require("@tailwindcss/forms")],
};
