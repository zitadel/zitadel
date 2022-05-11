const COLORS = [
  {
    500: "#ef4444",
    200: "#fecaca",
    300: "#fca5a5",
    600: "#dc2626",
    700: "#b91c1c",
    900: "#7f1d1d",
  },
  {
    500: "#f97316",
    200: "#fed7aa",
    300: "#fdba74",
    600: "#ea580c",
    700: "#c2410c",
    900: "#7c2d12",
  },
  {
    500: "#f59e0b",
    200: "#fde68a",
    300: "#fcd34d",
    600: "#d97706",
    700: "#b45309",
    900: "#78350f",
  },
  {
    500: "#eab308",
    200: "#fef08a",
    300: "#fde047",
    600: "#ca8a04",
    700: "#a16207",
    900: "#713f12",
  },
  {
    500: "#84cc16",
    200: "#d9f99d",
    300: "#bef264",
    600: "#65a30d",
    700: "#4d7c0f",
    900: "#365314",
  },
  {
    500: "#22c55e",
    200: "#bbf7d0",
    300: "#86efac",
    600: "#16a34a",
    700: "#15803d",
    900: "#14532d",
  },
  {
    500: "#10b981",
    200: "#a7f3d0",
    300: "#6ee7b7",
    600: "#059669",
    700: "#047857",
    900: "#064e3b",
  },
  {
    500: "#14b8a6",
    200: "#99f6e4",
    300: "#5eead4",
    600: "#0d9488",
    700: "#0f766e",
    900: "#134e4a",
  },
  {
    500: "#06b6d4",
    200: "#a5f3fc",
    300: "#67e8f9",
    600: "#0891b2",
    700: "#0e7490",
    900: "#164e63",
  },
  {
    500: "#0ea5e9",
    200: "#bae6fd",
    300: "#7dd3fc",
    600: "#0284c7",
    700: "#0369a1",
    900: "#0c4a6e",
  },
  {
    500: "#3b82f6",
    200: "#bfdbfe",
    300: "#93c5fd",
    600: "#2563eb",
    700: "#1d4ed8",
    900: "#1e3a8a",
  },
  {
    500: "#6366f1",
    200: "#c7d2fe",
    300: "#a5b4fc",
    600: "#4f46e5",
    700: "#4338ca",
    900: "#312e81",
  },
  {
    500: "#8b5cf6",
    200: "#ddd6fe",
    300: "#c4b5fd",
    600: "#7c3aed",
    700: "#6d28d9",
    900: "#4c1d95",
  },
  {
    500: "#a855f7",
    200: "#e9d5ff",
    300: "#d8b4fe",
    600: "#9333ea",
    700: "#7e22ce",
    900: "#581c87",
  },
  {
    500: "#d946ef",
    200: "#f5d0fe",
    300: "#f0abfc",
    600: "#c026d3",
    700: "#a21caf",
    900: "#701a75",
  },
  {
    500: "#ec4899",
    200: "#fbcfe8",
    300: "#f9a8d4",
    600: "#db2777",
    700: "#be185d",
    900: "#831843",
  },
  {
    500: "#f43f5e",
    200: "#fecdd3",
    300: "#fda4af",
    600: "#e11d48",
    700: "#be123c",
    900: "#881337",
  },
];

function hashCode(str, seed = 0) {
  let h1 = 0xdeadbeef ^ seed,
    h2 = 0x41c6ce57 ^ seed;
  for (let i = 0, ch; i < str.length; i++) {
    ch = str.charCodeAt(i);
    h1 = Math.imul(h1 ^ ch, 2654435761);
    h2 = Math.imul(h2 ^ ch, 1597334677);
  }
  h1 =
    Math.imul(h1 ^ (h1 >>> 16), 2246822507) ^
    Math.imul(h2 ^ (h2 >>> 13), 3266489909);
  h2 =
    Math.imul(h2 ^ (h2 >>> 16), 2246822507) ^
    Math.imul(h1 ^ (h1 >>> 13), 3266489909);
  return 4294967296 * (2097151 & h2) + (h1 >>> 0);
}

function getColor(value) {
  let hash = 0;

  if (value.length === 0) {
    return COLORS[hash];
  }

  hash = hashCode(value);
  return COLORS[hash % COLORS.length];
}

const avatars = document.getElementsByClassName("lgn-avatar");
for (let i = 0; i < avatars.length; i++) {
  const displayName = avatars[i].getAttribute("loginname");
  if (displayName) {
    const username = displayName.split("@")[0];
    let separator = "_";
    if (username.includes("-")) {
      separator = "-";
    }
    if (username.includes(".")) {
      separator = ".";
    }
    const split = username.split(separator);
    const initials = split[0].charAt(0) + (split[1] ? split[1].charAt(0) : "");
    avatars[i].getElementsByClassName("initials")[0].innerHTML = initials;

    const colorPalette = this.getColor(displayName);

    const isDark =
      document.documentElement.classList.contains("lgn-dark-theme");

    const backgroundShade = isDark ? 900 : 300;
    const foregroundShade = isDark ? 200 : 900;

    avatars[i].style.background = colorPalette[backgroundShade];
    avatars[i].style.color = colorPalette[foregroundShade];
  }
}
