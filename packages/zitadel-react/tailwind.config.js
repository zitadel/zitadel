const colors = require("tailwindcss/colors");

/** @type {import('tailwindcss').Config} */
module.exports = {
  prefix: "ui-",
  darkMode: "class",
  content: [`src/**/*.{js,ts,jsx,tsx}`],
  theme: {
    extend: {
      colors: {
        brandblue: colors.blue[500],
        brandred: colors.red[500],
        divider: {
          light: "rgba(135,149,161,.2)",
          dark: "rgba(135,149,161,.2)",
        },
      },
    },
  },
  plugins: [],
};
