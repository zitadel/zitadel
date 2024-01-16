import { render, screen } from "@testing-library/react";
import { SignInWithGoogle } from "./SignInWithGoogle";

describe("<SignInWithGoogle />", () => {
  it("renders without crashing", () => {
    const { container } = render(<SignInWithGoogle />);
    expect(container.firstChild).toBeDefined();
  });

  it("displays the default text", () => {
    render(<SignInWithGoogle />);
    const signInText = screen.getByText(/Sign in with Google/i);
    expect(signInText).toBeInTheDocument();
  });

  it("displays the given text", () => {
    render(<SignInWithGoogle name={"Google"} />);
    const signInText = screen.getByText(/Google/i);
    expect(signInText).toBeInTheDocument();
  });
});
