"use client";

import { createPortal } from "react-dom";
import { ReactNode, useEffect, useRef, useState } from "react";

/**
 * Client component that renders sanitized HTML and mounts React components
 * into placeholder elements via portals.
 *
 * Why manual innerHTML instead of dangerouslySetInnerHTML:
 * dangerouslySetInnerHTML causes React to re-create DOM nodes on re-render,
 * which invalidates the element references used by createPortal.
 * By setting innerHTML manually via a ref, the DOM is not managed by React
 * and survives across re-renders.
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

    // Set HTML manually — React does NOT manage this content,
    // so the DOM nodes persist across re-renders.
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
      {/* Server-renders the HTML for SEO; client replaces via useEffect */}
      <div
        ref={containerRef}
        suppressHydrationWarning
        dangerouslySetInnerHTML={{ __html: html }}
      />
      {ready &&
        slotRefs.current.map(([name, el]) =>
          slots[name] ? createPortal(slots[name], el, name) : null,
        )}
    </>
  );
}
