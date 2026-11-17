import { cleanup, fireEvent, render } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";
import { RegisterPrompt } from "./register-prompt";

const push = vi.fn();

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push }),
}));

vi.mock("next-intl", () => ({
  useTranslations: () => (key: string) => key,
}));

describe("RegisterPrompt", () => {
  beforeEach(() => vi.clearAllMocks());
  afterEach(cleanup);

  test("asks first time visitors to register", () => {
    const { getByTestId } = render(<RegisterPrompt />);

    expect(getByTestId("first-time-here").textContent).toBe("firstTimeHere");
    expect(getByTestId("register-button").textContent).toBe("register");
  });

  test("opens the registration with the organization and the request", () => {
    const { getByTestId } = render(<RegisterPrompt organization="org-1" requestId="oidc_1" />);

    fireEvent.click(getByTestId("register-button"));

    expect(push).toHaveBeenCalledWith("/register?organization=org-1&requestId=oidc_1");
  });
});
