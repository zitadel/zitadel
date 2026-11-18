import { cleanup, render } from "@testing-library/react";
import { afterEach, describe, expect, test, vi } from "vitest";
import { CopyToClipboard } from "./copy-to-clipboard";

vi.mock("next-intl", () => ({
  useTranslations: () => (key: string) => key,
}));

describe("CopyToClipboard", () => {
  afterEach(cleanup);

  test("should give the icon-only button an accessible name", () => {
    const { getByRole } = render(<CopyToClipboard value="SECRET" />);
    expect(getByRole("button", { name: "copy" })).toBeInTheDocument();
  });
});
