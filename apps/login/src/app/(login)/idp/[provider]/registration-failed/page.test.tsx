import Link from "next/link";
import { isValidElement, ReactNode } from "react";
import { describe, expect, test, vi } from "vitest";
import Page from "./page";

vi.mock("@/lib/zitadel", () => ({
  getBrandingSettings: vi.fn(async () => ({})),
  getDefaultOrg: vi.fn(async () => null),
}));

vi.mock("@/lib/service-url", () => ({
  getServiceConfig: vi.fn(() => ({ serviceConfig: { baseUrl: "https://example.com" } })),
}));

vi.mock("next-intl/server", () => ({
  getTranslations: vi.fn(async () => (key: string) => key),
}));

vi.mock("next/headers", () => ({
  headers: vi.fn(() => new Headers()),
}));

function findLinks(node: ReactNode): { href: string }[] {
  if (Array.isArray(node)) return node.flatMap(findLinks);
  if (!isValidElement<{ children?: ReactNode; href: string }>(node)) return [];
  return [...(node.type === Link ? [node.props] : []), ...findLinks(node.props.children)];
}

const render = (searchParams: Record<string, string>) => Page({ searchParams: Promise.resolve(searchParams) });

describe("registration-failed page", () => {
  test("renders no back link for an external postErrorRedirectUrl", async () => {
    const links = findLinks(await render({ postErrorRedirectUrl: "https://evil.example" }));
    expect(links).toEqual([]);
  });

  test("renders no back link without postErrorRedirectUrl", async () => {
    expect(findLinks(await render({}))).toEqual([]);
  });

  test("renders the back link for /loginname", async () => {
    const links = findLinks(await render({ postErrorRedirectUrl: "/loginname" }));
    expect(links.map((l) => l.href)).toEqual(["/loginname"]);
  });
});
