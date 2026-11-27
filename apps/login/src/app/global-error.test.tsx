import { cleanup, fireEvent, render } from "@testing-library/react";
import { afterEach, describe, expect, test, vi } from "vitest";
import GlobalError from "./global-error";

// next-intl is deliberately not mocked and no NextIntlClientProvider is rendered:
// global-error replaces the root layout, so it must render without either.
describe("GlobalError", () => {
  afterEach(cleanup);

  test("shows the error and a retry button without a next-intl provider", () => {
    const reset = vi.fn();
    // an <html> element cannot be rendered into the default <div> container
    const { getByText, getByRole } = render(<GlobalError error={new Error("Error in input stream")} reset={reset} />, {
      container: document,
    });

    expect(getByText("Error in input stream")).toBeTruthy();

    fireEvent.click(getByRole("button", { name: "Try Again" }));
    expect(reset).toHaveBeenCalledTimes(1);
  });
});
