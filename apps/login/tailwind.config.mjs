import sharedConfig from "@zitadel/tailwind-config/tailwind.config.mjs";

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
      colors[type][theme][`contrast-${shade}`] =
        `var(--theme-${theme}-${type}-contrast-${shade})`;
      colors[type][theme][`secondary-${shade}`] =
        `var(--theme-${theme}-${type}-secondary-${shade})`;
    });
  });
});

/** @type {import('tailwindcss').Config} */
export default {
  presets: [sharedConfig],
  darkMode: "class",
  content: ["./src/**/*.{js,ts,jsx,tsx}"],
  future: {
    hoverOnlyWhenSupported: true,
  },
  theme: {
    extend: {
      colors: {
        ...colors,
        state: {
          success: {
            light: {
              background: "#cbf4c9",
              color: "#0e6245",
            },
            dark: {
              background: "#68cf8340",
              color: "#cbf4c9",
            },
          },
          error: {
            light: {
              background: "#ffc1c1",
              color: "#620e0e",
            },
            dark: {
              background: "#af455359",
              color: "#ffc1c1",
            },
          },
          neutral: {
            light: {
              background: "#e4e7e4",
              color: "#000000",
            },
            dark: {
              background: "#1a253c",
              color: "#ffffff",
            },
          },
          alert: {
            light: {
              background: "#fbbf24",
              color: "#92400e",
            },
            dark: {
              background: "#92400e50",
              color: "#fbbf24",
            },
          },
        },
      },
      animation: {
        shake: "shake .8s cubic-bezier(.36,.07,.19,.97) both;",
      },
      keyframes: {
        shake: {
          "10%, 90%": {
            transform: "translate3d(-1px, 0, 0)",
          },

          "20%, 80%": {
            transform: "translate3d(2px, 0, 0)",
          },

          "30%, 50%, 70%": {
            transform: "translate3d(-4px, 0, 0)",
          },

          "40%, 60%": {
            transform: "translate3d(4px, 0, 0)",
          },
        },
      },
    },
  },
  plugins: [require("@tailwindcss/forms")],
};
