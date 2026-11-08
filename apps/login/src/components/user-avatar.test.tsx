import { cleanup, render, screen } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";
import { UserAvatar } from "./user-avatar";

vi.mock("@/lib/zitadel", () => ({
  getUserByID: vi.fn(),
}));

vi.mock("@/lib/service-url", () => ({
  getServiceConfig: vi.fn(() => ({ serviceConfig: { baseUrl: "https://example.com" } })),
}));

vi.mock("next/headers", () => ({
  headers: vi.fn(() => new Headers()),
}));

vi.mock("@/components/avatar", () => ({ Avatar: vi.fn(() => null) }));

const loginName = "anna@example.com@studio.example.com";

describe("UserAvatar", () => {
  let getUserByID: any;

  beforeEach(async () => {
    vi.clearAllMocks();
    const zitadel = await import("@/lib/zitadel");
    getUserByID = vi.mocked(zitadel.getUserByID);
    getUserByID.mockResolvedValue({ user: { userId: "user-1", username: "anna@example.com" } });
  });

  afterEach(() => {
    cleanup();
  });

  test("shows the login name when the suffix is not hidden", async () => {
    render(await UserAvatar({ loginName, showDropdown: false, userId: "user-1", hideLoginNameSuffix: false }));

    expect(screen.getByText(loginName)).toBeDefined();
    expect(getUserByID).not.toHaveBeenCalled();
  });

  test("shows the username when the branding hides the login name suffix", async () => {
    render(await UserAvatar({ loginName, showDropdown: false, userId: "user-1", hideLoginNameSuffix: true }));

    expect(screen.getByText("anna@example.com")).toBeDefined();
    expect(screen.queryByText(loginName)).toBeNull();
    expect(getUserByID).toHaveBeenCalledWith(expect.objectContaining({ userId: "user-1" }));
  });

  test("shows the login name when the suffix is hidden but the user is unknown", async () => {
    render(await UserAvatar({ loginName, showDropdown: false, hideLoginNameSuffix: true }));

    expect(screen.getByText(loginName)).toBeDefined();
    expect(getUserByID).not.toHaveBeenCalled();
  });

  test("shows the login name when the user cannot be read", async () => {
    getUserByID.mockRejectedValue(new Error("permission denied"));

    render(await UserAvatar({ loginName, showDropdown: false, userId: "user-1", hideLoginNameSuffix: true }));

    expect(screen.getByText(loginName)).toBeDefined();
  });
});
