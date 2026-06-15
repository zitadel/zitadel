"use client";

import { createPortal } from "react-dom";
import { ReactNode, useEffect, useRef, useState } from "react";
import { sanitizeLiquidOutput } from "@/lib/sanitize-liquid";

/**
 * Known slot names that the portal system recognises.
 * Any `data-liquid-slot` value NOT in this set is silently ignored
 * so a malicious template cannot create arbitrary portal mount points.
 */
const ALLOWED_SLOT_NAMES = new Set(["theme_switcher", "language_switcher"]);

/**
 * Client component that renders sanitized HTML and mounts React components
 * into placeholder elements via portals.
 *
 * How it works:
 * 1. The container div is initially rendered empty by React.
 * 2. On mount, `useEffect` re-sanitizes the html (defense-in-depth) and
 *    sets `container.innerHTML`.  Because React never "owns" these DOM
 *    nodes (no `dangerouslySetInnerHTML`), they persist across re‑renders.
 * 3. `querySelectorAll("[data-liquid-slot]")` finds the placeholder
 *    elements, filters against ALLOWED_SLOT_NAMES, and `createPortal`
 *    mounts the corresponding React components into them.
 *
 * Security:
 * - The `html` prop MUST have been sanitized by `sanitizeLiquidOutput`
 *   before reaching this component (done in `splitAtContent`).
 * - As defense-in-depth this component re-sanitizes the html before
 *   injecting it via innerHTML.
 * - Only known slot names are accepted; unknown `data-liquid-slot`
 *   values are ignored.
 *
 * Why NOT dangerouslySetInnerHTML:
 * React treats content set via `dangerouslySetInnerHTML` as its own.
 * When the component re‑renders (e.g. after `setReady(true)`), React
 * re‑applies `dangerouslySetInnerHTML`, which destroys the DOM nodes
 * that `createPortal` targets — the portals end up attached to
 * orphaned, off‑document elements and nothing appears on screen.
 */
export function LiquidSlotRenderer({
  html,
  slots,
}: {
  html: string;
  slots: Record<string, ReactNode>;
}) {
  const containerRef = useRef<HTMLDivElement>(null);
  const slotRefs = useRef<Array<[string, Element]>>([]);
  const [ready, setReady] = useState(false);

  useEffect(() => {
    const container = containerRef.current;
    if (!container) return;

    // Defense-in-depth: re-sanitize even though the caller already did.
    // This ensures safety if the call-site ever changes or a new caller
    // is introduced without proper sanitization.
    container.innerHTML = sanitizeLiquidOutput(html);

    // Find slot placeholder elements — only accept known slot names.
    const found: Array<[string, Element]> = [];
    container.querySelectorAll("[data-liquid-slot]").forEach((el) => {
      const name = el.getAttribute("data-liquid-slot");
      if (name && ALLOWED_SLOT_NAMES.has(name)) {
        found.push([name, el]);
      }
    });

    slotRefs.current = found;
    setReady(found.length > 0);

    return () => {
      slotRefs.current = [];
      setReady(false);
    };
  }, [html]);

  return (
    <>
      {/* Empty div — useEffect fills it with innerHTML on the client.
          suppressHydrationWarning is set because the server renders this
          div empty while the client fills it via useEffect. */}
      <div ref={containerRef} suppressHydrationWarning />
      {ready &&
        slotRefs.current.map(([name, el]) =>
          slots[name] ? createPortal(slots[name], el, name) : null,
        )}
    </>
  );
}
