const colors = require("tailwindcss/colors");

/** @type {import('tailwindcss').Config} */
module.exports = {
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
      // https://vercel.com/design/color
      fontSize: {
        "12px": "12px",
        "14px": "14px",
      },
      colors: {
        gray: colors.zinc,
        divider: {
          dark: "rgba(135,149,161,.2)",
          light: "rgba(135,149,161,.2)",
        },
        input: {
          light: {
            label: "#000000c7",
            background: "#00000004",
            border: "#1a191954",
            hoverborder: "1a1b1b",
          },
          dark: {
            label: "#ffffffc7",
            background: "#00000020",
            border: "#f9f7f775",
            hoverborder: "#e0e0e0",
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
      },
      backgroundImage: ({ theme }) => ({
        "dark-vc-border-gradient": `radial-gradient(at left top, ${theme(
          "colors.gray.800"
        )}, 50px, ${theme("colors.gray.800")} 50%)`,
        "vc-border-gradient": `radial-gradient(at left top, ${theme(
          "colors.gray.200"
        )}, 50px, ${theme("colors.gray.300")} 50%)`,
      }),
      keyframes: ({ theme }) => ({
        rerender: {
          "0%": {
            ["border-color"]: theme("colors.pink.500"),
          },
          "40%": {
            ["border-color"]: theme("colors.pink.500"),
          },
        },
        highlight: {
          "0%": {
            background: theme("colors.pink.500"),
            color: theme("colors.white"),
          },
          "40%": {
            background: theme("colors.pink.500"),
            color: theme("colors.white"),
          },
        },
        shimmer: {
          "100%": {
            transform: "translateX(100%)",
          },
        },
        translateXReset: {
          "100%": {
            transform: "translateX(0)",
          },
        },
        fadeToTransparent: {
          "0%": {
            opacity: 1,
          },
          "40%": {
            opacity: 1,
          },
          "100%": {
            opacity: 0,
          },
        },
      }),
    },
  },
  plugins: [require("@tailwindcss/forms")],
};
