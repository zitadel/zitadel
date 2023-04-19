const colors = require("tailwindcss/colors");

/** @type {import('tailwindcss').Config} */
module.exports = {
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
      // https://vercel.com/design/color
      fontSize: {
        "12px": "12px",
        "14px": "14px",
      },
      colors: {
        gray: colors.zinc,
        "gray-1000": "rgb(17,17,19)",
        "gray-1100": "rgb(10,10,11)",
        vercel: {
          pink: "#FF0080",
          blue: "#0070F3",
          cyan: "#50E3C2",
          orange: "#F5A623",
          violet: "#7928CA",
        },
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
          light: {
            500: colors.white,
            600: colors.gray[50],
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
            ["border-color"]: theme("colors.vercel.pink"),
          },
          "40%": {
            ["border-color"]: theme("colors.vercel.pink"),
          },
        },
        highlight: {
          "0%": {
            background: theme("colors.vercel.pink"),
            color: theme("colors.white"),
          },
          "40%": {
            background: theme("colors.vercel.pink"),
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
