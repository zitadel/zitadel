"use client";

import { startTransition, useCallback, useEffect, useRef, useState, useTransition } from "react";

/** Upper bound for the wait: a stalled page stream must not keep the user on the form. */
export const SETTLE_TIMEOUT_MS = 3000;

type Waiter = { id: number; resolve: () => void };

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
 * How: the action runs in a transition. Once its result is there, a state update is queued
 * in a second transition (after an `await` the code is no longer inside the first one) and
 * the promise resolves in the effect of the commit that contains that update.
 *
 * This relies on React 19 entangling all pending async transitions: the update queued here
 * commits together with the page update Next.js queued for the action, not before it. React
 * intends to scope this per action. Re-check in Firefox or Safari on every React or Next.js
 * upgrade that the redirect after a password login shows no error page; the unit test cannot
 * detect a change of this behaviour.
 *
 * Use it for every server action that is followed by `handleServerActionResponse`, and for
 * a cookie-writing action that runs right before one:
 *
 *   const runAction = useSettledAction();
 *   const response = await runAction(() => sendPassword(command));
 *   handleServerActionResponse(response, router, setError);
 *
 * Do not use it for an action that is followed by something needing the user's activation
 * (`navigator.credentials.get`): the wait can outlast it.
 *
 * Limits:
 * - A rejected action rejects right away.
 * - If the component unmounts while waiting (the re-rendered page redirected or failed), the
 *   promise resolves at that point. Once the component is gone there is no commit left to
 *   wait for, so a result that arrives later resolves right away.
 * - The wait ends after SETTLE_TIMEOUT_MS at the latest, a possible error flash then being
 *   the smaller problem.
 * - The wait also applies when the result is an internal redirect (`router.push`), which the
 *   browser would not cancel. Waiting only for external redirects and SAML posts would save
 *   one render of the current page there.
 */
export function useSettledAction(): <T>(action: () => Promise<T>) => Promise<T> {
  const [, startActionTransition] = useTransition();
  const [settled, setSettled] = useState(0);
  const lastId = useRef(0);
  const waiting = useRef<Waiter[]>([]);
  const mounted = useRef(false);

  // resolves the callers whose update is part of what has been committed
  const release = useCallback((upTo: number) => {
    const due = waiting.current.filter((waiter) => waiter.id <= upTo);
    waiting.current = waiting.current.filter((waiter) => waiter.id > upTo);
    due.forEach((waiter) => waiter.resolve());
  }, []);

  // runs after the commit that contains the page update of the action
  useEffect(() => {
    release(settled);
  }, [settled, release]);

  // never leave a caller hanging when the component goes away
  useEffect(() => {
    mounted.current = true;
    return () => {
      mounted.current = false;
      release(Infinity);
    };
  }, [release]);

  return useCallback(
    <T>(action: () => Promise<T>) =>
      new Promise<T>((resolve, reject) => {
        startActionTransition(async () => {
          try {
            const result = await action();
            if (!mounted.current) {
              // an update of an unmounted component is never committed, so nothing would release
              resolve(result);
              return;
            }

            const id = ++lastId.current;
            const timeout = setTimeout(() => release(id), SETTLE_TIMEOUT_MS);
            waiting.current.push({
              id,
              resolve: () => {
                clearTimeout(timeout);
                resolve(result);
              },
            });
            // a new transition, see above for why it commits together with the page update
            startTransition(() => setSettled((current) => Math.max(current, id)));
          } catch (error) {
            reject(error);
          }
        });
      }),
    [startActionTransition, release],
  );
}
