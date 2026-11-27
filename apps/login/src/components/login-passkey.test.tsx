import { act, cleanup, fireEvent, render, screen, waitFor, within } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { StrictMode } from "react";
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";
import { LoginPasskey } from "./login-passkey";

// Mock next/navigation
const mockPush = vi.fn();
vi.mock("next/navigation", () => ({
  useRouter: () => ({
    push: mockPush,
  }),
}));

// Mock server actions
vi.mock("@/lib/actions/passkeys", () => ({
  sendPasskey: vi.fn(),
}));

vi.mock("@/lib/actions/session", () => ({
  updateOrCreateSession: vi.fn(),
}));

// Mock navigator.credentials
const mockCredentialsGet = vi.fn();
Object.defineProperty(global.navigator, "credentials", {
  value: {
    get: mockCredentialsGet,
  },
  writable: true,
});

describe("LoginPasskey Component", () => {
  let mockSendPasskey: any;
  let mockUpdateSession: any;

  const messages = {
    error: {
      sessionResolution: {
        restart: "Restart sign-in",
      },
    },
    passkey: {
      verify: {
        title: "Authenticate with a passkey",
        description: "Your device will ask for your fingerprint",
        usePassword: "Use password",
        submit: "Continue",
        errors: {
          couldNotRequestChallenge: "Could not request passkey challenge",
          couldNotVerifyPasskey: "Could not verify passkey",
          noResponseReceived: "Passkey verification failed - no response received",
          noRedirectProvided: "Passkey verification completed but no redirect was provided",
          couldNotRetrievePasskey: "An error occurred while retrieving passkey",
          verificationCancelled: "Passkey verification was cancelled",
          verificationFailed: "An error occurred during passkey verification",
        },
      },
    },
  };

  const renderWithIntl = (component: React.ReactElement) => {
    return render(
      <NextIntlClientProvider locale="en" messages={messages}>
        {component}
      </NextIntlClientProvider>,
    );
  };

  beforeEach(async () => {
    vi.clearAllMocks();
    mockPush.mockClear();
    mockCredentialsGet.mockClear();

    const { sendPasskey } = await import("@/lib/actions/passkeys");
    const { updateOrCreateSession } = await import("@/lib/actions/session");

    mockSendPasskey = vi.mocked(sendPasskey);
    mockUpdateSession = vi.mocked(updateOrCreateSession);
  });

  // without it the renders of earlier tests stay in the document and their elements match later queries
  afterEach(cleanup);

  describe("Initialization and Challenge Request", () => {
    test("should display the error of the challenge request", async () => {
      mockUpdateSession.mockResolvedValue({
        error: "Challenge failed",
      });

      renderWithIntl(<LoginPasskey loginName="test@example.com" altPassword={false} />);

      await waitFor(() => {
        expect(screen.getByText("Challenge failed")).toBeInTheDocument();
      });
      expect(screen.queryByText("Could not request passkey challenge")).toBeNull();
      expect(screen.queryByTestId("restart-sign-in")).toBeNull();
    });

    test("should display translated error when the challenge request rejects", async () => {
      mockUpdateSession.mockRejectedValue(new Error("network down"));

      renderWithIntl(<LoginPasskey loginName="test@example.com" altPassword={false} />);

      await waitFor(() => {
        expect(screen.getByText("Could not request passkey challenge")).toBeInTheDocument();
      });
      expect(screen.queryByTestId("restart-sign-in")).toBeNull();
    });

    // U2F and passkey steps cannot continue without their session, the user has to start over
    test("should offer to restart the sign-in when the session of the step is lost", async () => {
      mockUpdateSession.mockResolvedValue({
        error: "The login session for this account is no longer available.",
        sessionLost: true,
      });

      renderWithIntl(
        <LoginPasskey
          loginName="test@example.com"
          organization="org-1"
          requestId="oidc_123"
          altPassword={false}
          login={false}
        />,
      );

      await waitFor(() => {
        expect(screen.getByText("The login session for this account is no longer available.")).toBeInTheDocument();
      });
      expect(screen.getByTestId("restart-sign-in")).toHaveAttribute(
        "href",
        "/loginname?loginName=test%40example.com&organization=org-1&requestId=oidc_123",
      );
    });

    test("should display translated error when no public key is returned", async () => {
      mockUpdateSession.mockResolvedValue({
        challenges: {
          webAuthN: {
            publicKeyCredentialRequestOptions: {},
          },
        },
      });

      renderWithIntl(<LoginPasskey loginName="test@example.com" altPassword={false} />);

      await waitFor(() => {
        expect(screen.getByText("Could not request passkey challenge")).toBeInTheDocument();
      });
    });

    test("should trigger passkey prompt when public key is available", async () => {
      const mockPublicKey = {
        challenge: new Uint8Array([1, 2, 3]),
        allowCredentials: [
          {
            id: new Uint8Array([4, 5, 6]),
            type: "public-key",
          },
        ],
      };

      mockUpdateSession.mockResolvedValue({
        challenges: {
          webAuthN: {
            publicKeyCredentialRequestOptions: {
              publicKey: mockPublicKey,
            },
          },
        },
      });

      mockCredentialsGet.mockResolvedValue({
        id: "credential-id",
        rawId: new ArrayBuffer(8),
        type: "public-key",
        response: {
          authenticatorData: new ArrayBuffer(8),
          clientDataJSON: new ArrayBuffer(8),
          signature: new ArrayBuffer(8),
          userHandle: new ArrayBuffer(8),
        },
      });

      mockSendPasskey.mockResolvedValue({
        redirect: "/success",
      });

      renderWithIntl(<LoginPasskey loginName="test@example.com" altPassword={false} />);

      await waitFor(() => {
        expect(mockCredentialsGet).toHaveBeenCalled();
      });
    });
  });

  describe("Passkey Verification", () => {
    const setupSuccessfulChallenge = () => {
      const mockPublicKey = {
        challenge: new Uint8Array([1, 2, 3]),
        allowCredentials: [
          {
            id: new Uint8Array([4, 5, 6]),
            type: "public-key",
          },
        ],
      };

      mockUpdateSession.mockResolvedValue({
        challenges: {
          webAuthN: {
            publicKeyCredentialRequestOptions: {
              publicKey: mockPublicKey,
            },
          },
        },
      });

      return mockPublicKey;
    };

    test("should redirect on successful verification", async () => {
      setupSuccessfulChallenge();

      mockCredentialsGet.mockResolvedValue({
        id: "credential-id",
        rawId: new ArrayBuffer(8),
        type: "public-key",
        response: {
          authenticatorData: new ArrayBuffer(8),
          clientDataJSON: new ArrayBuffer(8),
          signature: new ArrayBuffer(8),
          userHandle: new ArrayBuffer(8),
        },
      });

      mockSendPasskey.mockResolvedValue({
        redirect: "/success",
      });

      renderWithIntl(<LoginPasskey loginName="test@example.com" altPassword={false} />);

      await waitFor(() => {
        expect(mockPush).toHaveBeenCalledWith("/success");
      });
    });

    test("should display error when sendPasskey returns error", async () => {
      setupSuccessfulChallenge();

      mockCredentialsGet.mockResolvedValue({
        id: "credential-id",
        rawId: new ArrayBuffer(8),
        type: "public-key",
        response: {
          authenticatorData: new ArrayBuffer(8),
          clientDataJSON: new ArrayBuffer(8),
          signature: new ArrayBuffer(8),
          userHandle: new ArrayBuffer(8),
        },
      });

      mockSendPasskey.mockResolvedValue({
        error: "Verification failed",
      });

      renderWithIntl(<LoginPasskey loginName="test@example.com" altPassword={false} />);

      await waitFor(() => {
        expect(screen.getByText("Verification failed")).toBeInTheDocument();
      });
      expect(screen.queryByTestId("restart-sign-in")).toBeNull();
    });

    test("should offer to restart the sign-in when sendPasskey reports the session as lost", async () => {
      setupSuccessfulChallenge();

      mockCredentialsGet.mockResolvedValue({
        id: "credential-id",
        rawId: new ArrayBuffer(8),
        type: "public-key",
        response: {
          authenticatorData: new ArrayBuffer(8),
          clientDataJSON: new ArrayBuffer(8),
          signature: new ArrayBuffer(8),
          userHandle: new ArrayBuffer(8),
        },
      });

      mockSendPasskey.mockResolvedValue({
        error: "This login session has ended or expired.",
        sessionLost: true,
      });

      renderWithIntl(<LoginPasskey loginName="test@example.com" requestId="oidc_123" altPassword={false} />);

      await waitFor(() => {
        expect(screen.getByText("This login session has ended or expired.")).toBeInTheDocument();
      });
      expect(screen.getByTestId("restart-sign-in")).toHaveAttribute(
        "href",
        "/loginname?loginName=test%40example.com&requestId=oidc_123",
      );
    });

    test("should display error when sendPasskey returns undefined", async () => {
      setupSuccessfulChallenge();

      mockCredentialsGet.mockResolvedValue({
        id: "credential-id",
        rawId: new ArrayBuffer(8),
        type: "public-key",
        response: {
          authenticatorData: new ArrayBuffer(8),
          clientDataJSON: new ArrayBuffer(8),
          signature: new ArrayBuffer(8),
          userHandle: new ArrayBuffer(8),
        },
      });

      mockSendPasskey.mockResolvedValue(undefined);

      renderWithIntl(<LoginPasskey loginName="test@example.com" altPassword={false} />);

      await waitFor(() => {
        expect(screen.getByText("Passkey verification failed - no response received")).toBeInTheDocument();
      });
    });

    test("should display error when sendPasskey returns object without redirect", async () => {
      setupSuccessfulChallenge();

      mockCredentialsGet.mockResolvedValue({
        id: "credential-id",
        rawId: new ArrayBuffer(8),
        type: "public-key",
        response: {
          authenticatorData: new ArrayBuffer(8),
          clientDataJSON: new ArrayBuffer(8),
          signature: new ArrayBuffer(8),
          userHandle: new ArrayBuffer(8),
        },
      });

      mockSendPasskey.mockResolvedValue({});

      renderWithIntl(<LoginPasskey loginName="test@example.com" altPassword={false} />);

      await waitFor(() => {
        expect(screen.getByText("Passkey verification completed but no redirect was provided")).toBeInTheDocument();
      });
    });

    test("should display error when credential retrieval returns null", async () => {
      setupSuccessfulChallenge();
      mockCredentialsGet.mockResolvedValue(null);

      renderWithIntl(<LoginPasskey loginName="test@example.com" altPassword={false} />);

      await waitFor(() => {
        expect(screen.getByText("An error occurred while retrieving passkey")).toBeInTheDocument();
      });
    });
  });

  describe("Error Handling for Passkey Cancellation", () => {
    test("should display cancellation message when user cancels passkey", async () => {
      const mockPublicKey = {
        challenge: new Uint8Array([1, 2, 3]),
        allowCredentials: [
          {
            id: new Uint8Array([4, 5, 6]),
            type: "public-key",
          },
        ],
      };

      mockUpdateSession.mockResolvedValue({
        challenges: {
          webAuthN: {
            publicKeyCredentialRequestOptions: {
              publicKey: mockPublicKey,
            },
          },
        },
      });

      const notAllowedError = new Error("User cancelled");
      (notAllowedError as any).name = "NotAllowedError";
      mockCredentialsGet.mockRejectedValue(notAllowedError);

      renderWithIntl(<LoginPasskey loginName="test@example.com" altPassword={false} />);

      await waitFor(() => {
        expect(screen.getByText("Passkey verification was cancelled")).toBeInTheDocument();
      });
    });

    test("should display generic error for other credential errors", async () => {
      const mockPublicKey = {
        challenge: new Uint8Array([1, 2, 3]),
        allowCredentials: [
          {
            id: new Uint8Array([4, 5, 6]),
            type: "public-key",
          },
        ],
      };

      mockUpdateSession.mockResolvedValue({
        challenges: {
          webAuthN: {
            publicKeyCredentialRequestOptions: {
              publicKey: mockPublicKey,
            },
          },
        },
      });

      const genericError = new Error("Unknown error");
      mockCredentialsGet.mockRejectedValue(genericError);

      renderWithIntl(<LoginPasskey loginName="test@example.com" altPassword={false} />);

      await waitFor(() => {
        expect(screen.getByText("An error occurred during passkey verification")).toBeInTheDocument();
      });
    });
  });

  describe("Props Handling", () => {
    test("should pass sessionId to server actions when provided", async () => {
      mockUpdateSession.mockResolvedValue({
        error: "Test error",
      });

      renderWithIntl(<LoginPasskey sessionId="session-123" altPassword={false} />);

      await waitFor(() => {
        expect(mockUpdateSession).toHaveBeenCalledWith(
          expect.objectContaining({
            sessionId: "session-123",
          }),
        );
      });
    });

    test("should pass loginName to server actions when provided", async () => {
      mockUpdateSession.mockResolvedValue({
        error: "Test error",
      });

      renderWithIntl(<LoginPasskey loginName="test@example.com" altPassword={false} />);

      await waitFor(() => {
        expect(mockUpdateSession).toHaveBeenCalledWith(
          expect.objectContaining({
            loginName: "test@example.com",
          }),
        );
      });
    });

    test("should pass organization to server actions when provided", async () => {
      mockUpdateSession.mockResolvedValue({
        error: "Test error",
      });

      renderWithIntl(<LoginPasskey loginName="test@example.com" organization="org-123" altPassword={false} />);

      await waitFor(() => {
        expect(mockUpdateSession).toHaveBeenCalledWith(
          expect.objectContaining({
            organization: "org-123",
          }),
        );
      });
    });

    test("should pass requestId to server actions when provided", async () => {
      mockUpdateSession.mockResolvedValue({
        error: "Test error",
      });

      renderWithIntl(<LoginPasskey loginName="test@example.com" requestId="request-123" altPassword={false} />);

      await waitFor(() => {
        expect(mockUpdateSession).toHaveBeenCalledWith(
          expect.objectContaining({
            requestId: "request-123",
          }),
        );
      });
    });
  });

  describe("Single passkey ceremony", () => {
    const challengeResponse = () => ({
      challenges: {
        webAuthN: {
          publicKeyCredentialRequestOptions: {
            publicKey: {
              challenge: new Uint8Array([1, 2, 3]),
              allowCredentials: [{ id: new Uint8Array([4, 5, 6]), type: "public-key" }],
            },
          },
        },
      },
    });

    const assertion = {
      id: "credential-id",
      rawId: new ArrayBuffer(8),
      type: "public-key",
      response: {
        authenticatorData: new ArrayBuffer(8),
        clientDataJSON: new ArrayBuffer(8),
        signature: new ArrayBuffer(8),
        userHandle: new ArrayBuffer(8),
      },
    };

    // Keeps navigator.credentials.get pending, as while the browser's passkey prompt is open.
    const holdPrompt = () => {
      const prompt: {
        signal?: AbortSignal;
        resolve: (value: unknown) => void;
        reject: (reason: unknown) => void;
      } = { resolve: () => {}, reject: () => {} };
      mockCredentialsGet.mockImplementation(
        (options: CredentialRequestOptions) =>
          new Promise((resolve, reject) => {
            prompt.signal = options.signal ?? undefined;
            prompt.resolve = resolve;
            prompt.reject = reject;
          }),
      );
      return prompt;
    };

    const notAllowed = () => Object.assign(new Error("cancelled"), { name: "NotAllowedError" });

    test("keeps Continue disabled and starts no second ceremony while the automatic prompt is open", async () => {
      mockUpdateSession.mockResolvedValue(challengeResponse());
      const prompt = holdPrompt();

      const { container } = renderWithIntl(<LoginPasskey loginName="test@example.com" altPassword={false} />);
      const continueButton = within(container).getByTestId("submit-button");

      await waitFor(() => {
        expect(mockCredentialsGet).toHaveBeenCalledTimes(1);
      });
      expect(continueButton).toBeDisabled();

      fireEvent.click(continueButton);
      fireEvent.click(continueButton);

      expect(mockUpdateSession).toHaveBeenCalledTimes(1);
      expect(mockCredentialsGet).toHaveBeenCalledTimes(1);

      prompt.reject(notAllowed());

      await waitFor(() => {
        expect(continueButton).not.toBeDisabled();
      });
      expect(within(container).getByText("Passkey verification was cancelled")).toBeInTheDocument();
    });

    test("keeps Continue disabled while a ceremony started by Continue is open", async () => {
      mockUpdateSession.mockResolvedValue(challengeResponse());
      mockCredentialsGet.mockRejectedValueOnce(notAllowed());

      const { container } = renderWithIntl(<LoginPasskey loginName="test@example.com" altPassword={false} />);
      const continueButton = within(container).getByTestId("submit-button");

      await waitFor(() => {
        expect(continueButton).not.toBeDisabled();
      });

      const prompt = holdPrompt();
      fireEvent.click(continueButton);

      await waitFor(() => {
        expect(mockCredentialsGet).toHaveBeenCalledTimes(2);
      });
      expect(continueButton).toBeDisabled();

      fireEvent.click(continueButton);
      expect(mockUpdateSession).toHaveBeenCalledTimes(2);

      prompt.reject(notAllowed());

      await waitFor(() => {
        expect(continueButton).not.toBeDisabled();
      });
    });

    test("aborts the open prompt when switching to password", async () => {
      mockUpdateSession.mockResolvedValue(challengeResponse());
      const prompt = holdPrompt();

      const { container } = renderWithIntl(<LoginPasskey loginName="test@example.com" altPassword={true} />);

      await waitFor(() => {
        expect(mockCredentialsGet).toHaveBeenCalledTimes(1);
      });
      expect(prompt.signal?.aborted).toBe(false);

      fireEvent.click(within(container).getByTestId("password-button"));

      expect(prompt.signal?.aborted).toBe(true);
      expect(mockPush).toHaveBeenCalledWith(expect.stringMatching(/^\/password\?/));

      // A late result of the aborted prompt is ignored.
      await act(async () => {
        prompt.resolve(assertion);
      });
      expect(mockSendPasskey).not.toHaveBeenCalled();
      expect(within(container).queryByRole("alert")).not.toBeInTheDocument();
    });

    test("aborts the open prompt when the page is left another way", async () => {
      mockUpdateSession.mockResolvedValue(challengeResponse());
      const prompt = holdPrompt();

      const { unmount } = renderWithIntl(<LoginPasskey loginName="test@example.com" altPassword={false} />);

      await waitFor(() => {
        expect(mockCredentialsGet).toHaveBeenCalledTimes(1);
      });

      // Back, or any other navigation, unmounts the component.
      unmount();
      expect(prompt.signal?.aborted).toBe(true);

      await act(async () => {
        prompt.resolve(assertion);
      });
      expect(mockSendPasskey).not.toHaveBeenCalled();
    });

    test("opens a single prompt under strict mode", async () => {
      mockUpdateSession.mockResolvedValue(challengeResponse());
      const prompt = holdPrompt();

      render(
        <StrictMode>
          <NextIntlClientProvider locale="en" messages={messages}>
            <LoginPasskey loginName="test@example.com" altPassword={false} />
          </NextIntlClientProvider>
        </StrictMode>,
      );

      await waitFor(() => {
        expect(mockCredentialsGet).toHaveBeenCalledTimes(1);
      });
      expect(prompt.signal?.aborted).toBe(false);
    });

    test("ignores the verification result when switching to password during verification", async () => {
      mockUpdateSession.mockResolvedValue(challengeResponse());
      mockCredentialsGet.mockResolvedValue(assertion);
      let finishVerification: (value: unknown) => void = () => {};
      mockSendPasskey.mockImplementation(
        () =>
          new Promise((resolve) => {
            finishVerification = resolve;
          }),
      );

      const { container } = renderWithIntl(<LoginPasskey loginName="test@example.com" altPassword={true} />);

      await waitFor(() => {
        expect(mockSendPasskey).toHaveBeenCalledTimes(1);
      });

      fireEvent.click(within(container).getByTestId("password-button"));
      await act(async () => {
        finishVerification({ redirect: "/success" });
      });

      expect(mockPush).toHaveBeenCalledTimes(1);
      expect(mockPush).toHaveBeenCalledWith(expect.stringMatching(/^\/password\?/));
    });

    test("re-enables Continue when the page is restored from the back/forward cache", async () => {
      mockUpdateSession.mockResolvedValue(challengeResponse());
      mockCredentialsGet.mockResolvedValue(assertion);
      mockSendPasskey.mockResolvedValue({ redirect: "/success" });

      const { container } = renderWithIntl(<LoginPasskey loginName="test@example.com" altPassword={false} />);
      const continueButton = within(container).getByTestId("submit-button");

      await waitFor(() => {
        expect(mockPush).toHaveBeenCalledWith("/success");
      });
      expect(continueButton).toBeDisabled();

      const pageShow = new Event("pageshow");
      Object.defineProperty(pageShow, "persisted", { value: true });
      act(() => {
        window.dispatchEvent(pageShow);
      });

      expect(continueButton).not.toBeDisabled();
    });

    test("re-enables Continue some time after an external redirect", async () => {
      vi.useFakeTimers({ shouldAdvanceTime: true });
      const location = window.location;
      try {
        mockUpdateSession.mockResolvedValue(challengeResponse());
        mockCredentialsGet.mockResolvedValue(assertion);
        // A native app's custom-scheme callback opens the app and leaves this page in place.
        mockSendPasskey.mockResolvedValue({ redirect: "myapp://callback" });
        const hrefSetter = vi.fn();
        Object.defineProperty(window, "location", {
          configurable: true,
          value: {
            ...location,
            set href(value: string) {
              hrefSetter(value);
            },
          },
        });

        const { container } = renderWithIntl(<LoginPasskey loginName="test@example.com" altPassword={false} />);
        const continueButton = within(container).getByTestId("submit-button");

        await waitFor(() => {
          expect(hrefSetter).toHaveBeenCalledWith("myapp://callback");
        });
        expect(continueButton).toBeDisabled();

        act(() => {
          vi.advanceTimersByTime(5000);
        });
        expect(continueButton).not.toBeDisabled();
      } finally {
        Object.defineProperty(window, "location", { configurable: true, value: location });
        vi.useRealTimers();
      }
    });

    test("keeps Continue disabled while redirecting after a successful verification", async () => {
      mockUpdateSession.mockResolvedValue(challengeResponse());
      mockCredentialsGet.mockResolvedValue(assertion);
      mockSendPasskey.mockResolvedValue({ redirect: "/success" });

      const { container } = renderWithIntl(<LoginPasskey loginName="test@example.com" altPassword={false} />);
      const continueButton = within(container).getByTestId("submit-button");

      await waitFor(() => {
        expect(mockPush).toHaveBeenCalledWith("/success");
      });
      expect(continueButton).toBeDisabled();

      fireEvent.click(continueButton);
      expect(mockUpdateSession).toHaveBeenCalledTimes(1);
    });

    test("re-enables Continue when verification fails", async () => {
      mockUpdateSession.mockResolvedValue(challengeResponse());
      mockCredentialsGet.mockResolvedValue(assertion);
      mockSendPasskey.mockResolvedValue({ error: "Verification failed" });

      const { container } = renderWithIntl(<LoginPasskey loginName="test@example.com" altPassword={false} />);

      await waitFor(() => {
        expect(within(container).getByText("Verification failed")).toBeInTheDocument();
      });
      expect(within(container).getByTestId("submit-button")).not.toBeDisabled();
    });
  });

  describe("useEffect Initialization Guard", () => {
    test("should only initialize once even if component re-renders", async () => {
      mockUpdateSession.mockResolvedValue({
        error: "Test error",
      });

      const { rerender } = renderWithIntl(<LoginPasskey loginName="test@example.com" altPassword={false} />);

      await waitFor(() => {
        expect(mockUpdateSession).toHaveBeenCalledTimes(1);
      });

      // Rerender the component
      rerender(
        <NextIntlClientProvider locale="en" messages={messages}>
          <LoginPasskey loginName="test@example.com" altPassword={false} />
        </NextIntlClientProvider>,
      );

      // Should still only be called once
      expect(mockUpdateSession).toHaveBeenCalledTimes(1);
    });
  });
});
