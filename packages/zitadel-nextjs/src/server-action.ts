/**
 * Higher-order function that wraps a server action with authentication checks.
 * Placeholder — to be implemented with Next.js server action primitives.
 */
export function protectedAction<TArgs extends unknown[], TResult>(
  action: (...args: TArgs) => Promise<TResult>,
): (...args: TArgs) => Promise<TResult> {
  return async (...args: TArgs) => {
    // TODO: verify session before executing action
    return action(...args);
  };
}
