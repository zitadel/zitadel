// @vitest-environment node
import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { describe, expect, test } from "vitest";

// A key that is missing renders as the raw key (e.g. `register.errors.registerNotAllowed`) and
// logs MISSING_MESSAGE, so every locale must hold every key the code uses.

const loginRoot = fileURLToPath(new URL("../..", import.meta.url));
const localesDir = path.join(loginRoot, "locales");

type Messages = { [key: string]: string | Messages };
type KeyKinds = Record<string, "leaf" | "object">;

/** Every key of a message tree, nested objects included, e.g. `device.request.title`. */
function keysOf(messages: Messages, prefix = "", keys: KeyKinds = {}): KeyKinds {
  for (const [key, value] of Object.entries(messages)) {
    const fullKey = prefix ? `${prefix}.${key}` : key;
    if (typeof value === "object") {
      keys[fullKey] = "object";
      keysOf(value, fullKey, keys);
    } else {
      keys[fullKey] = "leaf";
    }
  }
  return keys;
}

const leavesOf = (keys: KeyKinds) => Object.keys(keys).filter((key) => keys[key] === "leaf");

const locales = fs
  .readdirSync(localesDir)
  .filter((file) => file.endsWith(".json"))
  .map((file) => file.slice(0, -".json".length));

const keysByLocale: Record<string, KeyKinds> = {};
for (const locale of locales) {
  keysByLocale[locale] = keysOf(JSON.parse(fs.readFileSync(path.join(localesDir, `${locale}.json`), "utf8")));
}
const englishKeys = keysByLocale.en;

function sourceFiles(dir: string): string[] {
  return fs.readdirSync(dir, { withFileTypes: true }).flatMap((entry) => {
    const file = path.join(dir, entry.name);
    if (entry.isDirectory()) {
      return sourceFiles(file);
    }
    return /\.tsx?$/.test(entry.name) && !/\.test\.tsx?$/.test(entry.name) ? [file] : [];
  });
}

function matches(pattern: RegExp, source: string): RegExpExecArray[] {
  const found: RegExpExecArray[] = [];
  let match: RegExpExecArray | null;
  while ((match = pattern.exec(source)) !== null) {
    found.push(match);
  }
  return found;
}

/**
 * The keys the source uses with a literal key: `t("key")` on a translator from
 * `getTranslations`/`useTranslations`, and `<Translated i18nKey="key" namespace="ns" />`.
 * Keys built at runtime are not covered.
 */
function usedKeys(): { key: string; where: string }[] {
  const used: { key: string; where: string }[] = [];

  for (const file of sourceFiles(path.join(loginRoot, "src"))) {
    const source = fs.readFileSync(file, "utf8");
    const where = (index: number) => `${path.relative(loginRoot, file)}:${source.slice(0, index).split("\n").length}`;

    // one variable can hold translators of different namespaces in different functions
    const namespacesByTranslator: Record<string, string[]> = {};
    const translator =
      /(?:const|let)\s+(\w+)\s*=\s*(?:await\s+)?(?:getTranslations|useTranslations)\(\s*(?:"([^"]*)"|\{[^}]*namespace:\s*"([^"]*)"[^}]*\}|\{[^}]*\}|)\s*\)/g;
    for (const match of matches(translator, source)) {
      const namespaces = (namespacesByTranslator[match[1]] ??= []);
      namespaces.push(match[2] ?? match[3] ?? "");
    }

    for (const [name, namespaces] of Object.entries(namespacesByTranslator)) {
      const call = new RegExp(`\\b${name}(?:\\.rich|\\.markup|\\.raw)?\\(\\s*["'\`]([^"'\`$]+)["'\`]`, "g");
      for (const match of matches(call, source)) {
        const candidates = namespaces.map((namespace) => (namespace ? `${namespace}.${match[1]}` : match[1]));
        used.push({
          key: candidates.find((key) => key in englishKeys) ?? candidates.join(" | "),
          where: where(match.index),
        });
      }
    }

    // `[^>]` also matches line breaks, so attributes on several lines are found too
    for (const match of matches(/<Translated\b([^>]*?)\/?>/g, source)) {
      const key = /i18nKey=(?:"([^"]+)"|\{"([^"]+)"\})/.exec(match[1]);
      if (!key) {
        continue; // built at runtime
      }
      const namespace = /namespace=(?:"([^"]+)"|\{"([^"]+)"\})/.exec(match[1]);
      const name = key[1] ?? key[2];
      const ns = namespace ? (namespace[1] ?? namespace[2]) : "";
      used.push({ key: ns ? `${ns}.${name}` : name, where: where(match.index) });
    }
  }

  return used;
}

describe("login messages", () => {
  test("finds the keys the source uses", () => {
    // guards the scan itself: if the patterns stop matching, the next test would pass vacuously
    expect(usedKeys().length).toBeGreaterThan(200);
  });

  test("every key the source uses exists in English", () => {
    const missing = usedKeys()
      .filter(({ key }) => !(key in englishKeys))
      .map(({ key, where }) => `${key} (${where})`);

    expect(missing).toEqual([]);
  });

  test.each(locales.filter((locale) => locale !== "en"))("%s has exactly the English keys", (locale) => {
    const keys = keysByLocale[locale];

    expect(leavesOf(englishKeys).filter((key) => keys[key] !== "leaf")).toEqual([]);
    expect(leavesOf(keys).filter((key) => englishKeys[key] !== "leaf")).toEqual([]);
  });
});
