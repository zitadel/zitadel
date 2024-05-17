import { afterEach, describe, expect, test } from "vitest";

import { cleanup, render, screen } from "@testing-library/react";
import { SignInWithGitlab } from "./SignInWithGitlab";

afterEach(cleanup);

describe("<SignInWithGitlab />", () => {
  test("renders without crashing", () => {
    const { container } = render(<SignInWithGitlab />);
    expect(container.firstChild).toBeDefined();
  });

  test("displays the default text", () => {
    render(<SignInWithGitlab />);
    const signInText = screen.getByText(/Sign in with Gitlab/i);
    expect(signInText).toBeInTheDocument();
  });

  test("displays the given text", () => {
    render(<SignInWithGitlab name={"Gitlab"} />);
    const signInText = screen.getByText(/Gitlab/i);
    expect(signInText).toBeInTheDocument();
  });
});
