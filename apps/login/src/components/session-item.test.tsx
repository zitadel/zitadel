import { cleanup, fireEvent, render, waitFor } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";
import { SessionItem } from "./session-item";

const { push } = vi.hoisted(() => ({ push: vi.fn() }));

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push }),
}));

vi.mock("next-intl", () => ({
  useTranslations: () => (key: string) => key,
  useLocale: () => "en",
}));

vi.mock("next-themes", () => ({
  useTheme: () => ({ theme: "light" }),
}));

vi.mock("@radix-ui/react-tooltip", () => ({
  Root: ({ children }: any) => children,
  Trigger: ({ children }: any) => children,
  Portal: () => null,
  Content: () => null,
  Arrow: () => null,
}));

vi.mock("@/lib/server/loginname", () => ({
  sendLoginname: vi.fn(),
}));

vi.mock("@/lib/server/session", () => ({
  clearSession: vi.fn(),
  continueWithSession: vi.fn(),
}));

describe("SessionItem", () => {
  let sendLoginname: any;

  beforeEach(async () => {
    vi.clearAllMocks();
    const loginname = await import("@/lib/server/loginname");
    sendLoginname = vi.mocked(loginname.sendLoginname);
  });

  afterEach(cleanup);

  const unverifiedSession = {
    factors: {
      user: {
        loginName: "user@example.com",
        displayName: "Example User",
        organizationId: "123",
      },
    },
  } as any;

  test("should show a translated error when continuing an unverified session fails", async () => {
    sendLoginname.mockRejectedValue(new Error("internal error"));

    const { getByRole, getByText } = render(
      <SessionItem session={unverifiedSession} reload={() => {}} requestId="request_123" />,
    );

    fireEvent.click(getByRole("button"));

    await waitFor(() => expect(getByText("internalError")).toBeInTheDocument());
    expect(sendLoginname).toHaveBeenCalledWith({
      loginName: "user@example.com",
      organization: "123",
      requestId: "request_123",
    });
    expect(push).not.toHaveBeenCalled();
  });
});
