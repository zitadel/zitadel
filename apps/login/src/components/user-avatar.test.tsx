import { cleanup, render, screen } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";
import { UserAvatar } from "./user-avatar";

vi.mock("next-themes", () => ({ useTheme: () => ({ resolvedTheme: "light" }) }));
afterEach(cleanup);

describe("UserAvatar", () => {
  it("shows the profile image and preserves the account switch link", () => {
    render(
      <UserAvatar
        loginName="alex@example.com"
        displayName="Alex Rivera"
        imageUrl="https://idp.example.com/assets/v1/photo"
        showDropdown
        searchParams={{ requestId: "oidc-1", organization: "org-1" }}
      />,
    );
    expect(screen.getByRole("img", { name: "avatar" })).toHaveAttribute("src", "https://idp.example.com/assets/v1/photo");
    expect(screen.getByRole("link")).toHaveAttribute("href", "/accounts?organization=org-1&requestId=oidc-1");
  });

  it("shows initials without a profile image", () => {
    render(<UserAvatar loginName="alex@example.com" displayName="Alex Rivera" showDropdown={false} />);
    expect(screen.queryByRole("img")).not.toBeInTheDocument();
    expect(screen.getByText("AR")).toBeInTheDocument();
  });
});
