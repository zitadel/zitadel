import { cleanup, fireEvent, render, waitFor } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";
import { ConsentScreen } from "./consent";

const { push } = vi.hoisted(() => ({ push: vi.fn() }));

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push }),
}));

vi.mock("next-intl", () => ({
  useTranslations: () => (key: string) => key,
}));

vi.mock("@/lib/server/device", () => ({
  approveDeviceAuthorization: vi.fn(),
  denyDeviceAuthorization: vi.fn(),
}));

describe("ConsentScreen", () => {
  let approveDeviceAuthorization: any;
  let denyDeviceAuthorization: any;

  beforeEach(async () => {
    vi.clearAllMocks();
    const device = await import("@/lib/server/device");
    approveDeviceAuthorization = vi.mocked(device.approveDeviceAuthorization);
    denyDeviceAuthorization = vi.mocked(device.denyDeviceAuthorization);

    approveDeviceAuthorization.mockResolvedValue({ redirect: "/signedin?requestId=device_ABCD-1234" });
    denyDeviceAuthorization.mockResolvedValue({ details: {} });
  });

  afterEach(cleanup);

  test("shows the user code so it can be compared with the device", () => {
    const { getByTestId } = render(<ConsentScreen userCode="ABCD-1234" sessionId="session-1" />);

    expect(getByTestId("user-code")).toHaveTextContent("ABCD-1234");
  });

  test("approves the device authorization and follows the returned redirect", async () => {
    const { getByTestId } = render(
      <ConsentScreen userCode="ABCD-1234" sessionId="session-1" organization="org-1" appName="My Device" />,
    );

    fireEvent.click(getByTestId("submit-button"));

    await waitFor(() => expect(approveDeviceAuthorization).toHaveBeenCalledTimes(1));
    expect(approveDeviceAuthorization).toHaveBeenCalledWith({
      userCode: "ABCD-1234",
      sessionId: "session-1",
      organization: "org-1",
    });
    await waitFor(() => expect(push).toHaveBeenCalledWith("/signedin?requestId=device_ABCD-1234"));
  });

  test("denies the device authorization and returns to the code entry page", async () => {
    const { getByTestId } = render(<ConsentScreen userCode="ABCD-1234" sessionId="session-1" organization="org-1" />);

    fireEvent.click(getByTestId("deny-button"));

    await waitFor(() =>
      expect(denyDeviceAuthorization).toHaveBeenCalledWith({
        userCode: "ABCD-1234",
        sessionId: "session-1",
        organization: "org-1",
      }),
    );
    await waitFor(() => expect(push).toHaveBeenCalledWith("/device"));
    expect(approveDeviceAuthorization).not.toHaveBeenCalled();
  });
});
