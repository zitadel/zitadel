import {
  lowerCaseValidator,
  numberValidator,
  symbolValidator,
  upperCaseValidator,
} from "@/helpers/validators";
import { PasswordComplexitySettings } from "@zitadel/proto/zitadel/settings/v2/password_settings_pb";

type Props = {
  passwordComplexitySettings: PasswordComplexitySettings;
  password: string;
  equals: boolean;
};

const check = (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    fill="none"
    viewBox="0 0 24 24"
    strokeWidth={1.5}
    stroke="currentColor"
    className="las la-check mr-2 h-6 w-6 text-lg text-green-500 dark:text-green-500"
    role="img"
  >
    <title>Matches</title>
    <path
      strokeLinecap="round"
      strokeLinejoin="round"
      d="M4.5 12.75l6 6 9-13.5"
    />
  </svg>
);
const cross = (
  <svg
    className="las la-times mr-2 h-6 w-6 text-lg text-warn-light-500 dark:text-warn-dark-500"
    xmlns="http://www.w3.org/2000/svg"
    fill="none"
    viewBox="0 0 24 24"
    strokeWidth={1.5}
    stroke="currentColor"
    role="img"
  >
    <title>Doesn&apos;t match</title>
    <path
      strokeLinecap="round"
      strokeLinejoin="round"
      d="M6 18L18 6M6 6l12 12"
    />
  </svg>
);
const desc =
  "text-14px leading-4 text-input-light-label dark:text-input-dark-label";

export function PasswordComplexity({
  passwordComplexitySettings,
  password,
  equals,
}: Props) {
  const hasMinLength = password?.length >= passwordComplexitySettings.minLength;
  const hasSymbol = symbolValidator(password);
  const hasNumber = numberValidator(password);
  const hasUppercase = upperCaseValidator(password);
  const hasLowercase = lowerCaseValidator(password);

  return (
    <div className="mb-4 grid grid-cols-2 gap-x-8 gap-y-2">
      {passwordComplexitySettings.minLength != undefined ? (
        <div className="flex flex-row items-center" data-testid="length-check">
          {hasMinLength ? check : cross}
          <span className={desc}>
            Password length {passwordComplexitySettings.minLength.toString()}
          </span>
        </div>
      ) : (
        <span />
      )}
      <div className="flex flex-row items-center" data-testid="symbol-check">
        {hasSymbol ? check : cross}
        <span className={desc}>has Symbol</span>
      </div>
      <div className="flex flex-row items-center" data-testid="number-check">
        {hasNumber ? check : cross}
        <span className={desc}>has Number</span>
      </div>
      <div className="flex flex-row items-center" data-testid="uppercase-check">
        {hasUppercase ? check : cross}
        <span className={desc}>has uppercase</span>
      </div>
      <div className="flex flex-row items-center" data-testid="lowercase-check">
        {hasLowercase ? check : cross}
        <span className={desc}>has lowercase</span>
      </div>
      <div className="flex flex-row items-center" data-testid="equal-check">
        {equals ? check : cross}
        <span className={desc}>equals</span>
      </div>
    </div>
  );
}
