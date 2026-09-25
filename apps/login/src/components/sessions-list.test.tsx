import { Provider as TooltipProvider } from "@radix-ui/react-tooltip";
import { cleanup, render, screen, within } from "@testing-library/react";
import { create } from "@zitadel/client";
import { SessionSchema } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { afterEach, describe, expect, it, vi } from "vitest";
import { SessionsList } from "./sessions-list";

vi.mock("@/lib/server/session", () => ({ clearSession: vi.fn(), continueWithSession: vi.fn() }));
vi.mock("@/lib/server/loginname", () => ({ sendLoginname: vi.fn() }));
vi.mock("next/navigation", () => ({ useRouter: () => ({ push: vi.fn() }) }));
vi.mock("next-intl", () => ({ useLocale: () => "en", useTranslations: () => (key: string) => key }));
vi.mock("next-themes", () => ({ useTheme: () => ({ resolvedTheme: "light" }) }));
vi.mock("./translated", () => ({ Translated: () => null }));
afterEach(cleanup);

describe("SessionsList avatars", () => {
  it("matches photos to users and retains initials for users without photos", () => {
    const sessions = [
      create(SessionSchema, {
        id: "session-1",
        factors: { user: { id: "alex", loginName: "alex@example.com", displayName: "Alex Rivera" } },
      }),
      create(SessionSchema, {
        id: "session-2",
        factors: { user: { id: "sam", loginName: "sam@example.com", displayName: "Sam Chen" } },
      }),
    ];
    render(
      <TooltipProvider>
        <SessionsList sessions={sessions} avatarUrls={{ alex: "https://idp.example.com/assets/v1/alex" }} />
      </TooltipProvider>,
    );
    const alex = screen.getByRole("button", { name: /Alex Rivera/ });
    const sam = screen.getByRole("button", { name: /Sam Chen/ });
    expect(within(alex).getByRole("img")).toHaveAttribute("src", "https://idp.example.com/assets/v1/alex");
    expect(within(sam).queryByRole("img")).not.toBeInTheDocument();
    expect(within(sam).getByText("SC")).toBeInTheDocument();
  });
});
