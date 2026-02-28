import type { ChangeEventHandler, FormEvent } from "react";

export interface UsernamePasswordFormProps {
  loginName: string;
  password: string;
  onLoginNameChange: ChangeEventHandler<HTMLInputElement>;
  onPasswordChange: ChangeEventHandler<HTMLInputElement>;
  onSubmit: (event: FormEvent<HTMLFormElement>) => void | Promise<void>;
  isLoading?: boolean;
  error?: string | null;
}

/**
 * Presentational username/password form with submit/loading/error states.
 */
export function UsernamePasswordForm({
  loginName,
  password,
  onLoginNameChange,
  onPasswordChange,
  onSubmit,
  isLoading = false,
  error,
}: UsernamePasswordFormProps) {
  return (
    <>
      <form onSubmit={onSubmit}>
        <h3>1) Create session</h3>
        <p>
          <label>
            Login name
            <br />
            <input
              autoComplete="username"
              name="loginName"
              onChange={onLoginNameChange}
              placeholder="mini@mouse.com"
              required
              value={loginName}
            />
          </label>
        </p>
        <p>
          <label>
            Password
            <br />
            <input
              autoComplete="current-password"
              name="password"
              onChange={onPasswordChange}
              required
              type="password"
              value={password}
            />
          </label>
        </p>
        <p>
          <button disabled={isLoading} type="submit">
            {isLoading ? "Creating session..." : "Create session"}
          </button>
        </p>
      </form>

      {error ? (
        <p className="ztdl-status-error" role="alert">
          <strong>Session error:</strong> {error}
        </p>
      ) : null}
    </>
  );
}
