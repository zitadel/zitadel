import { act, cleanup, render } from "@testing-library/react";
import { startTransition, Suspense, use, useState } from "react";
import { afterEach, describe, expect, test, vi } from "vitest";
import { useSettledAction } from "./settled-action";

type Run = ReturnType<typeof useSettledAction>;

function deferred<T>() {
  let resolve!: (value: T) => void;
  const promise = new Promise<T>((r) => (resolve = r));
  return { promise, resolve };
}

/**
 * `page` stands in for the page Next.js re-renders after a server action wrote a cookie:
 * a transition update that suspends until the rest of the action's response has arrived.
 */
function Harness({ expose }: { expose: (api: { run: Run; setPage: (page: Promise<string>) => void }) => void }) {
  const run = useSettledAction();
  const [page, setPage] = useState<Promise<string> | null>(null);
  expose({ run, setPage });
  return <p data-testid="page">{page ? use(page) : "initial"}</p>;
}

function setup() {
  let api!: { run: Run; setPage: (page: Promise<string>) => void };
  const utils = render(
    <Suspense fallback="loading">
      <Harness expose={(a) => (api = a)} />
    </Suspense>,
  );
  return { ...utils, api: () => api };
}

describe("useSettledAction", () => {
  afterEach(cleanup);

  test("resolves with the result of the action", async () => {
    const { api } = setup();

    // not awaited inside act(): the promise settles in an effect, which act() flushes on exit
    let result: unknown;
    await act(async () => {
      api()
        .run(async () => ({ redirect: "https://app.example.com/callback" }))
        .then((r) => (result = r));
    });

    expect(result).toEqual({ redirect: "https://app.example.com/callback" });
  });

  test("rejects when the action rejects", async () => {
    const { api } = setup();

    let error: unknown;
    await act(async () => {
      await api()
        .run(async () => {
          throw new Error("action failed");
        })
        .catch((e) => (error = e));
    });

    expect(error).toEqual(new Error("action failed"));
  });

  // Starts an action the way a cookie-writing server action behaves: the page update is queued
  // when the action is called, the result arrives while React is already waiting for the page.
  async function runActionWithPendingPage(api: () => { run: Run; setPage: (page: Promise<string>) => void }) {
    const page = deferred<string>();
    const result = deferred<string>();
    const settled = vi.fn();

    await act(async () => {
      api()
        .run(() => {
          // what Next.js does when the action is called: queue the page update in a transition
          startTransition(() => api().setPage(page.promise));
          return result.promise;
        })
        .then(settled);
    });

    await act(async () => {
      result.resolve("result");
    });

    return { page, settled };
  }

  test("resolves only after the page update the action triggered has been committed", async () => {
    const { api, getByTestId } = setup();
    const { page, settled } = await runActionWithPendingPage(api);

    // the action returned, but the rest of the re-rendered page has not arrived yet
    expect(settled).not.toHaveBeenCalled();
    expect(getByTestId("page").textContent).toBe("initial");

    await act(async () => {
      page.resolve("re-rendered");
    });

    expect(getByTestId("page").textContent).toBe("re-rendered");
    expect(settled).toHaveBeenCalledWith("result");
  });

  test("resolves when the component unmounts while waiting", async () => {
    const { api, unmount } = setup();
    const { settled } = await runActionWithPendingPage(api);

    expect(settled).not.toHaveBeenCalled();

    await act(async () => {
      unmount();
    });

    expect(settled).toHaveBeenCalledWith("result");
  });
});
