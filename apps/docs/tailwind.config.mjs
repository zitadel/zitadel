import typography from '@tailwindcss/typography';
import animate from 'tailwindcss-animate';
import { createRequire } from 'module';

const require = createRequire(import.meta.url);
const sharedConfig = require('@zitadel/theme/tailwind');

/** @type {import('tailwindcss').Config} */
export default {
  presets: [sharedConfig],
  darkMode: ["class", '[data-theme="dark"]'],
  content: [
    "./app/**/*.{js,ts,jsx,tsx}",
    "./components/**/*.{js,ts,jsx,tsx}",
    "./content/**/*.{md,mdx}",
    
    // FAIL-SAFE: Check both local and root node_modules for Fumadocs
    "./node_modules/fumadocs-ui/dist/**/*.js",
    "../../node_modules/fumadocs-ui/dist/**/*.js",
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ["var(--font-lato)", "sans-serif"],
      },
    },
  },
  plugins: [
    typography,
    animate,
  ],
};