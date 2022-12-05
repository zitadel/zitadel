const colors = require("tailwindcss/colors");

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./pages/**/*.{js,ts,jsx,tsx}",
    "./components/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      fontSize: {
        "12px": "12px",
        "14px": "14px",
      },
      colors: {
        gray: colors.zinc,
        primary: {
          light: {
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
        button: {
          light: {
            border: "#0000001f",
          },
          dark: {
            border: "#ffffff1f",
          },
        },
        border: {
          light: "rgba(135,149,161,.2)",
          dark: "rgba(135,149,161,.2)",
        },
      },
      lineHeight: {
        "14px": "14px",
        "14.5px": "14.5px",
        "36px": "36px",
        4: "1rem",
      },
    },
  },
  plugins: [],
};
