import { cleanup, fireEvent, render, waitFor } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";
import { OTPAddForm } from "./otp-add-form";

const push = vi.fn();

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push, back: vi.fn() }),
}));

vi.mock("next-intl", () => ({
  useTranslations: () => (key: string) => key,
}));

vi.mock("@/lib/actions/otp", () => ({
  addOTP: vi.fn(),
}));

describe("OTPAddForm", () => {
  let addOTP: any;

  beforeEach(async () => {
    vi.clearAllMocks();
    addOTP = vi.mocked((await import("@/lib/actions/otp")).addOTP);
  });

  afterEach(cleanup);

  test("does not add the factor until the user confirms", () => {
    render(<OTPAddForm method="email" sessionId="session-1" urlToContinue="/signedin" />);

    expect(addOTP).not.toHaveBeenCalled();
  });

  test("shows the success message and continue link after adding", async () => {
    addOTP.mockResolvedValue({ success: true });
    const { getByTestId, findByText } = render(
      <OTPAddForm method="email" sessionId="session-1" urlToContinue="/signedin?loginName=user" />,
    );

    fireEvent.click(getByTestId("add-button"));

    expect(await findByText("set.emailAdded")).toBeTruthy();
    expect(addOTP).toHaveBeenCalledWith({ method: "email", sessionId: "session-1" });
    expect(getByTestId("submit-button").closest("a")?.getAttribute("href")).toBe("/signedin?loginName=user");
    expect(push).not.toHaveBeenCalled();
  });

  test("continues to the code check directly with checkAfter", async () => {
    addOTP.mockResolvedValue({ success: true });
    const { getByTestId } = render(
      <OTPAddForm method="sms" sessionId="session-1" urlToContinue="/otp/sms?loginName=user" checkAfter />,
    );

    fireEvent.click(getByTestId("add-button"));

    await waitFor(() => expect(push).toHaveBeenCalledWith("/otp/sms?loginName=user"));
  });

  test("shows the error and no success message when adding fails", async () => {
    addOTP.mockResolvedValue({ error: "errors.couldNotAddSms" });
    const { getByTestId, findByText, queryByText } = render(
      <OTPAddForm method="sms" sessionId="session-1" urlToContinue="/signedin" />,
    );

    fireEvent.click(getByTestId("add-button"));

    expect(await findByText("errors.couldNotAddSms")).toBeTruthy();
    expect(queryByText("set.smsAdded")).toBeNull();
    expect(push).not.toHaveBeenCalled();
  });

  test("shows a translated error when the action fails unexpectedly", async () => {
    addOTP.mockRejectedValue(new Error("network"));
    const { getByTestId, findByText } = render(
      <OTPAddForm method="email" sessionId="session-1" urlToContinue="/signedin" />,
    );

    fireEvent.click(getByTestId("add-button"));

    expect(await findByText("errors.couldNotAddEmail")).toBeTruthy();
    expect(push).not.toHaveBeenCalled();
  });
});
