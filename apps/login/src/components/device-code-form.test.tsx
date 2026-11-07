import { cleanup, fireEvent, render, waitFor } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";
import { DeviceCodeForm } from "./device-code-form";

const { push } = vi.hoisted(() => ({ push: vi.fn() }));

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push }),
}));

vi.mock("next-intl", () => ({
  useTranslations: () => (key: string) => key,
}));

vi.mock("@/lib/server/oidc", () => ({
  getDeviceAuthorizationRequest: vi.fn(),
}));

describe("DeviceCodeForm", () => {
  let getDeviceAuthorizationRequest: any;

  beforeEach(async () => {
    vi.clearAllMocks();
    const oidc = await import("@/lib/server/oidc");
    getDeviceAuthorizationRequest = vi.mocked(oidc.getDeviceAuthorizationRequest);
  });

  afterEach(cleanup);

  async function submit(getByTestId: (id: string) => HTMLElement) {
    await waitFor(() => expect(getByTestId("submit-button")).not.toBeDisabled());
    fireEvent.click(getByTestId("submit-button"));
  }

  test("should autofocus the code input on mount", () => {
    const { getByTestId } = render(<DeviceCodeForm />);
    expect(getByTestId("code-text-input")).toHaveFocus();
  });

  test("should continue to the url returned by the server", async () => {
    getDeviceAuthorizationRequest.mockResolvedValue({
      deviceAuthorizationRequest: { id: "encrypted-device-code" },
      nextUrl: "/accounts?requestId=device_ABCD-1234",
    });

    const { getByTestId } = render(<DeviceCodeForm userCode="ABCD-1234" />);
    await submit(getByTestId);

    await waitFor(() => expect(push).toHaveBeenCalledWith("/accounts?requestId=device_ABCD-1234"));
    expect(getDeviceAuthorizationRequest).toHaveBeenCalledWith("ABCD-1234");
  });

  test("should show an error and not navigate when the server returns no next url", async () => {
    getDeviceAuthorizationRequest.mockResolvedValue({
      deviceAuthorizationRequest: { id: "encrypted-device-code" },
    });

    const { getByTestId } = render(<DeviceCodeForm userCode="ABCD-1234" />);
    await submit(getByTestId);

    await waitFor(() => expect(getByTestId("error")).toBeInTheDocument());
    expect(push).not.toHaveBeenCalled();
  });

  test("should show an error and not navigate when the code is unknown", async () => {
    getDeviceAuthorizationRequest.mockResolvedValue({ deviceAuthorizationRequest: undefined, nextUrl: undefined });

    const { getByTestId } = render(<DeviceCodeForm userCode="WXYZ-9876" />);
    await submit(getByTestId);

    await waitFor(() => expect(getByTestId("error")).toBeInTheDocument());
    expect(push).not.toHaveBeenCalled();
  });
});
