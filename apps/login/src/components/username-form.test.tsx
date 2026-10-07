import { cleanup, render } from "@testing-library/react";
import { afterEach, describe, expect, test, vi } from "vitest";
import { UsernameForm } from "./username-form";

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: vi.fn() }),
}));

vi.mock("next-intl", () => ({
  useTranslations: () => (key: string) => key,
}));

vi.mock("@/lib/server/loginname", () => ({
  sendLoginname: vi.fn(),
}));

describe("UsernameForm", () => {
  afterEach(cleanup);

  test("should autofocus the loginName input on mount", () => {
    const { getByTestId } = render(
      <UsernameForm loginName="" requestId={undefined} loginSettings={undefined} submit={false} allowRegister={false} />,
    );
    expect(getByTestId("username-text-input")).toHaveFocus();
  });

  test("shows the combined label when all methods are enabled", () => {
    const { getByText } = render(
      <UsernameForm
        loginName=""
        requestId={undefined}
        loginSettings={{ disableLoginWithEmail: false, disableLoginWithPhone: false } as any}
        submit={false}
        allowRegister={false}
      />,
    );
    expect(getByText("labels.usernameEmailOrPhone")).toBeInTheDocument();
  });

  test("shows username-only label when email and phone are disabled", () => {
    const { getByText } = render(
      <UsernameForm
        loginName=""
        requestId={undefined}
        loginSettings={{ disableLoginWithEmail: true, disableLoginWithPhone: true } as any}
        submit={false}
        allowRegister={false}
      />,
    );
    expect(getByText("labels.username")).toBeInTheDocument();
  });

  test("shows username-or-phone label when email is disabled", () => {
    const { getByText } = render(
      <UsernameForm
        loginName=""
        requestId={undefined}
        loginSettings={{ disableLoginWithEmail: true, disableLoginWithPhone: false } as any}
        submit={false}
        allowRegister={false}
      />,
    );
    expect(getByText("labels.usernameOrPhoneNumber")).toBeInTheDocument();
  });

  test("shows username-or-email label when phone is disabled", () => {
    const { getByText } = render(
      <UsernameForm
        loginName=""
        requestId={undefined}
        loginSettings={{ disableLoginWithEmail: false, disableLoginWithPhone: true } as any}
        submit={false}
        allowRegister={false}
      />,
    );
    expect(getByText("labels.usernameOrEmail")).toBeInTheDocument();
  });
});
