import { cleanup, render } from "@testing-library/react";
import { afterEach, describe, expect, test, vi } from "vitest";
import { DeviceCodeForm } from "./device-code-form";

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: vi.fn() }),
}));

vi.mock("next-intl", () => ({
  useTranslations: () => (key: string) => key,
}));

vi.mock("@/lib/server/oidc", () => ({
  getDeviceAuthorizationRequest: vi.fn(),
}));

describe("DeviceCodeForm", () => {
  afterEach(cleanup);

  test("should autofocus the code input on mount", () => {
    const { getByTestId } = render(<DeviceCodeForm />);
    expect(getByTestId("code-text-input")).toHaveFocus();
  });
});
