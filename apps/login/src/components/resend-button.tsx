import { Translated } from "./translated";

type Props = {
  namespace: string;
  /** label while the button can be used, e.g. `verify.resendCode` */
  i18nKey: string;
  /** label while the cooldown runs; receives `{ seconds }`, e.g. `verify.resendCodeIn` */
  countdownI18nKey: string;
  /** the state of `useResendCooldown` */
  cooldown: { remaining: number; isCoolingDown: boolean };
  disabled?: boolean;
  onClick: () => void;
};

/** The "Resend code" button of the code pages: disabled and counting down while the resend cooldown runs. */
export function ResendButton({ namespace, i18nKey, countdownI18nKey, cooldown, disabled, onClick }: Props) {
  return (
    <button
      disabled={disabled || cooldown.isCoolingDown}
      type="button"
      className="text-primary-light-500 hover:text-primary-light-400 dark:text-primary-dark-500 hover:dark:text-primary-dark-400 ml-4 cursor-pointer disabled:cursor-default disabled:text-gray-500 dark:disabled:text-gray-400"
      onClick={onClick}
      data-testid="resend-button"
    >
      {cooldown.isCoolingDown ? (
        <Translated i18nKey={countdownI18nKey} namespace={namespace} data={{ seconds: cooldown.remaining }} />
      ) : (
        <Translated i18nKey={i18nKey} namespace={namespace} />
      )}
    </button>
  );
}
