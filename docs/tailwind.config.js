/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: ["class", '[data-theme="dark"]'],
  content: ["./docs/**/*.{js,ts,jsx,tsx}", "./src/**/*.{js,jsx,ts,tsx}"],
  corePlugins: { preflight: false },
  theme: {
    fontFamily: {
      sans: ["Lato", "serif"],
    },
    extend: {
      colors: {
        zitadelpink: "#ff2069",
        primary: {
          dark: {
            100: "#afd1f2",
            200: "#7fb5ea",
            300: "#4192e0",
            400: "#2782dc",
            500: "#2073c4",
            600: "#1c64aa",
            700: "#17548f",
            800: "#134575",
            900: "#0f355b",
          },
        },
        background: {
          dark: {
            100: "#4a69aa",
            200: "#395183",
            300: "#243252",
            400: "#1a253c",
            500: "#111827",
            600: "#080b12",
            700: "#000000",
            800: "#000000",
            900: "#000000",
          },
        },
        input: {
          light: {
            label: "#000000c7",
            background: "#00000004",
            border: "rgba(26,25,25,.2196078431);",
          },
          dark: {
            label: "#ffffffc7",
            background: "#00000020",
            border: "rgba(249,247,247,.1450980392)",
          },
        },
      },
    },
  },
  plugins: [],
};
