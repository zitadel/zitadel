"use client";

import { startTransition, useCallback, useEffect, useRef, useState, useTransition } from "react";

/**
 * Returns a function that runs a server action and resolves with its result only after
 * the UI update the action triggered has been committed.
 *
 * Why: a server action that writes a cookie (the session cookie, for example) makes Next.js
 * re-render the current page and stream it in the same response as the action result. The
 * result arrives first. Leaving the page right away (`window.location.href`, a form post)
 * makes Firefox and Safari cancel the rest of that response, React fails the parts of the
 * page that are still missing, and an error page shows until the next page has loaded.
 *
 * The action runs inside a transition and a state update is queued in the same transition
 * once the result is there. React commits that update together with the re-rendered page,
 * so when the returned promise resolves nothing is left to cancel.
 *
 * Use it for every server action that is followed by `handleServerActionResponse`, and for
 * a cookie-writing action that runs before one:
 *
 *   const runAction = useSettledAction();
 *   const response = await runAction(() => sendPassword(command));
 *   handleServerActionResponse(response, router, setSamlData, setError);
 *
 * A rejected action rejects right away. If the component unmounts while waiting (the
 * re-rendered page redirected or failed), the promise resolves at that point. Once the
 * component is gone there is no commit left to wait for, so a result that arrives later
 * resolves right away.
 */
export function useSettledAction(): <T>(action: () => Promise<T>) => Promise<T> {
  const [, startActionTransition] = useTransition();
  const [settled, setSettled] = useState(0);
  const waiting = useRef<(() => void)[]>([]);
  const mounted = useRef(false);

  const flush = useCallback(() => {
    const callbacks = waiting.current;
    waiting.current = [];
    callbacks.forEach((callback) => callback());
  }, []);

  // runs after the commit that contains the page update of the action
  useEffect(() => {
    flush();
  }, [settled, flush]);

  // never leave a caller hanging when the component goes away
  useEffect(() => {
    mounted.current = true;
    return () => {
      mounted.current = false;
      flush();
    };
  }, [flush]);

  return useCallback(
    <T>(action: () => Promise<T>) =>
      new Promise<T>((resolve, reject) => {
        startActionTransition(async () => {
          try {
            const result = await action();
            if (!mounted.current) {
              // an update of an unmounted component is never committed, so nothing would flush
              resolve(result);
              return;
            }
            waiting.current.push(() => resolve(result));
            // still inside the action's transition, so this commits together with the page update
            startTransition(() => setSettled((count) => count + 1));
          } catch (error) {
            reject(error);
          }
        });
      }),
    [startActionTransition],
  );
}
