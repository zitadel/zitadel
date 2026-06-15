"use client";

import { createPortal } from "react-dom";
import { ReactNode, useEffect, useRef, useState } from "react";

/**
 * Client component that renders sanitized HTML and mounts React components
 * into placeholder elements via portals.
 *
 * How it works:
 * 1. The container div is initially rendered empty by React.
 * 2. On mount, `useEffect` sets `container.innerHTML` to inject the
 *    sanitized HTML.  Because React never "owns" these DOM nodes (no
 *    `dangerouslySetInnerHTML`), they persist across re‑renders.
 * 3. `querySelectorAll("[data-liquid-slot]")` finds the placeholder
 *    elements, and `createPortal` mounts the corresponding React
 *    components into them.
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

    // Inject HTML manually so React does NOT manage these DOM nodes.
    // They persist across re-renders, keeping portal targets valid.
    container.innerHTML = html;

    // Find slot placeholder elements
    const found: Array<[string, Element]> = [];
    container.querySelectorAll("[data-liquid-slot]").forEach((el) => {
      const name = el.getAttribute("data-liquid-slot");
      if (name) found.push([name, el]);
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
