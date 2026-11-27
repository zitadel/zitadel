import * as Tooltip from "@radix-ui/react-tooltip";
import { cleanup, fireEvent, render, screen, waitFor } from "@testing-library/react";
import { create } from "@zitadel/client";
import { SessionSchema } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";
import { SessionItem } from "./session-item";

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: vi.fn() }),
}));

vi.mock("next-intl", () => ({
  useLocale: () => "en",
  useTranslations: () => (key: string) => key,
}));

vi.mock("@/lib/actions/loginname", () => ({
  sendLoginname: vi.fn(),
}));

vi.mock("@/lib/server/session", () => ({
  clearSession: vi.fn(),
  continueWithSession: vi.fn(),
}));

vi.mock("@/lib/client-utils", () => ({
  handleServerActionResponse: vi.fn(),
}));

vi.mock("./avatar", () => ({ Avatar: vi.fn(() => null) }));
vi.mock("./translated", () => ({ Translated: vi.fn(() => null) }));

const loginName = "anna@example.com@studio.example.com";

// A session whose server side is gone (as after a logout): the account list still offers it.
const endedSession = create(SessionSchema, {
  id: "session-1",
  factors: { user: { loginName, organizationId: "org-1" } },
});

function renderItem(shownLoginName?: string) {
  return render(
    <Tooltip.Provider>
      <SessionItem session={endedSession} reload={vi.fn()} requestId="oidc_1" shownLoginName={shownLoginName} />
    </Tooltip.Provider>,
  );
}

describe("SessionItem", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  afterEach(() => {
    cleanup();
  });

  test("shows the login name it is given instead of the session's", () => {
    renderItem("anna@example.com");

    expect(screen.getByText("anna@example.com")).toBeDefined();
    expect(screen.queryByText(loginName)).toBeNull();
  });

  test("shows the session's login name when it is given none", () => {
    renderItem();

    expect(screen.getByText(loginName)).toBeDefined();
  });

  test("signs in with the session's own login name, not the one shown", async () => {
    const { sendLoginname } = await import("@/lib/actions/loginname");
    vi.mocked(sendLoginname).mockResolvedValue({ redirect: "/password" } as never);

    renderItem("anna@example.com");
    fireEvent.click(screen.getByRole("button"));

    await waitFor(() =>
      expect(sendLoginname).toHaveBeenCalledWith({ loginName, organization: "org-1", requestId: "oidc_1" }),
    );
  });
});
