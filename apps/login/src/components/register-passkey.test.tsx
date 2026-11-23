import { cleanup, fireEvent, render, screen, waitFor } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";
import { RegisterPasskey } from "./register-passkey";

const push = vi.fn();

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push }),
}));

vi.mock("next-intl", () => ({
  useTranslations: () => (key: string) => key,
}));

vi.mock("@/lib/actions/passkeys", () => ({
  registerPasskeyLink: vi.fn(),
  verifyPasskeyRegistration: vi.fn(),
}));

describe("RegisterPasskey", () => {
  let registerPasskeyLink: any;
  let verifyPasskeyRegistration: any;

  beforeEach(async () => {
    vi.clearAllMocks();
    const actions = await import("@/lib/actions/passkeys");
    registerPasskeyLink = vi.mocked(actions.registerPasskeyLink);
    verifyPasskeyRegistration = vi.mocked(actions.verifyPasskeyRegistration);

    registerPasskeyLink.mockResolvedValue({
      passkeyId: "passkey-1",
      publicKeyCredentialCreationOptions: { publicKey: { challenge: "Y2hhbGxlbmdl", user: { id: "dXNlcg" } } },
    });
    verifyPasskeyRegistration.mockResolvedValue({ loginName: "user@example.com" });

    const buffer = new Uint8Array([1, 2, 3]).buffer;
    Object.defineProperty(navigator, "credentials", {
      configurable: true,
      value: {
        create: vi.fn().mockResolvedValue({
          id: "cred-1",
          type: "public-key",
          rawId: buffer,
          response: { attestationObject: buffer, clientDataJSON: buffer },
        }),
      },
    });
  });

  afterEach(cleanup);

  test("uses the emailed code instead of a session in the browser", async () => {
    render(<RegisterPasskey sessionId="session-1" userId="user-1" code="code" codeId="code-id" isPrompt={false} />);

    await waitFor(() => expect(verifyPasskeyRegistration).toHaveBeenCalled());

    expect(registerPasskeyLink).toHaveBeenCalledWith({ userId: "user-1", code: "code", codeId: "code-id" });
    expect(verifyPasskeyRegistration).toHaveBeenCalledWith(
      expect.objectContaining({ passkeyId: "passkey-1", sessionId: undefined, userId: "user-1" }),
    );
    // the sign-in continues with the code's user, not with the session in the browser
    await waitFor(() => expect(push).toHaveBeenCalled());
    const url = new URLSearchParams(push.mock.calls[0][0].split("?")[1]);
    expect(url.get("sessionId")).toBeNull();
    expect(url.get("userId")).toBe("user-1");
    expect(url.get("loginName")).toBe("user@example.com");
  });

  test("uses the session when there is no code", async () => {
    render(<RegisterPasskey sessionId="session-1" isPrompt={false} />);

    const submit = screen.getByTestId("submit-button");
    await waitFor(() => expect(submit).not.toBeDisabled());
    fireEvent.click(submit);

    await waitFor(() => expect(verifyPasskeyRegistration).toHaveBeenCalled());

    expect(registerPasskeyLink).toHaveBeenCalledWith({ sessionId: "session-1" });
    expect(verifyPasskeyRegistration).toHaveBeenCalledWith(expect.objectContaining({ sessionId: "session-1" }));
    await waitFor(() => expect(push).toHaveBeenCalled());
    expect(new URLSearchParams(push.mock.calls[0][0].split("?")[1]).get("sessionId")).toBe("session-1");
  });

  test("shows the error and stays on the page when the verification is rejected", async () => {
    verifyPasskeyRegistration.mockResolvedValue({ error: "You have to verify your second factor first" });

    render(<RegisterPasskey sessionId="session-1" isPrompt={false} />);

    const submit = screen.getByTestId("submit-button");
    await waitFor(() => expect(submit).not.toBeDisabled());
    fireEvent.click(submit);

    await waitFor(() => expect(screen.getByText("You have to verify your second factor first")).toBeInTheDocument());
    expect(push).not.toHaveBeenCalled();
  });
});
