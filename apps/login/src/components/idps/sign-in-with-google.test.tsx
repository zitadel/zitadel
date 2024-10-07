import { afterEach, describe, expect, test } from "vitest";

import { cleanup, render, screen } from "@testing-library/react";
import { SignInWithGoogle } from "./sign-in-with-google";

afterEach(cleanup);

describe("<SignInWithGoogle />", () => {
  test("renders without crashing", () => {
    const { container } = render(<SignInWithGoogle />);
    expect(container.firstChild).toBeDefined();
  });

  test("displays the default text", () => {
    render(<SignInWithGoogle />);
    const signInText = screen.getByText(/Sign in with Google/i);
    expect(signInText).toBeInTheDocument();
  });

  test("displays the given text", () => {
    render(<SignInWithGoogle name={"Google"} />);
    const signInText = screen.getByText(/Google/i);
    expect(signInText).toBeInTheDocument();
  });
});
