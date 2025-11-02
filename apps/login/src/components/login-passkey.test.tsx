import { describe, expect, test, vi, beforeEach } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import { LoginPasskey } from "./login-passkey";
import { NextIntlClientProvider } from "next-intl";

// Mock next/navigation
const mockPush = vi.fn();
vi.mock("next/navigation", () => ({
  useRouter: () => ({
    push: mockPush,
  }),
}));

// Mock server actions
vi.mock("@/lib/server/passkeys", () => ({
  sendPasskey: vi.fn(),
}));

vi.mock("@/lib/server/session", () => ({
  updateSession: vi.fn(),
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

    const { sendPasskey } = await import("@/lib/server/passkeys");
    const { updateSession } = await import("@/lib/server/session");

    mockSendPasskey = vi.mocked(sendPasskey);
    mockUpdateSession = vi.mocked(updateSession);
  });

  describe("Initialization and Challenge Request", () => {
    test("should display error when challenge request fails", async () => {
      mockUpdateSession.mockResolvedValue({
        error: "Challenge failed",
      });

      renderWithIntl(<LoginPasskey loginName="test@example.com" altPassword={false} />);

      await waitFor(() => {
        expect(screen.getByText("Could not request passkey challenge")).toBeInTheDocument();
      });
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
