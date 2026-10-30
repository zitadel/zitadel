type Messages = Record<string, unknown>;

// keys that would change an object's prototype instead of adding a message
const UNSAFE_KEYS = new Set(["__proto__", "constructor", "prototype"]);

function isPlainObject(value: unknown): value is Messages {
  if (value === null || typeof value !== "object" || Array.isArray(value)) {
    return false;
  }
  const proto = Object.getPrototypeOf(value);
  return proto === Object.prototype || proto === null;
}

function mergeInto(target: Messages, source: Messages): Messages {
  for (const [key, value] of Object.entries(source)) {
    if (UNSAFE_KEYS.has(key)) {
      continue;
    }
    const existing = Object.hasOwn(target, key) ? target[key] : undefined;
    if (isPlainObject(value)) {
      target[key] = mergeInto(isPlainObject(existing) ? existing : {}, value);
    } else {
      target[key] = value;
    }
  }
  return target;
}

/**
 * Deep-merges message catalogs into a new object. Later sources win, nested objects are merged,
 * and any other value (strings, arrays) replaces the earlier one. Keys that could change an
 * object's prototype are skipped, because the last source comes from the hosted login translations.
 */
export function mergeMessages(...sources: Array<Messages | null | undefined>): Messages {
  const result: Messages = {};
  for (const source of sources) {
    if (isPlainObject(source)) {
      mergeInto(result, source);
    }
  }
  return result;
}
