import {
  lowerCaseValidator,
  numberValidator,
  symbolValidator,
  upperCaseValidator,
} from "@/helpers/validators";
import { PasswordComplexitySettings } from "@zitadel/proto/zitadel/settings/v2/password_settings_pb";
import { useTranslations } from "next-intl";
import { Translated } from "@/components/translated";

type Props = {
  passwordComplexitySettings: PasswordComplexitySettings;
  password: string;
  equals: boolean;
};

function CheckIcon({ title }: { title: string }) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      fill="none"
      viewBox="0 0 24 24"
      strokeWidth={1.5}
      stroke="currentColor"
      className="las la-check mr-2 h-6 w-6 flex-none text-lg text-green-500 dark:text-green-500"
      role="img"
    >
      <title>{title}</title>
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M4.5 12.75l6 6 9-13.5"
      />
    </svg>
  );
}

function CrossIcon({ title }: { title: string }) {
  return (
    <svg
      className="las la-times mr-2 h-6 w-6 flex-none text-lg text-warn-light-500 dark:text-warn-dark-500"
      xmlns="http://www.w3.org/2000/svg"
      fill="none"
      viewBox="0 0 24 24"
      strokeWidth={1.5}
      stroke="currentColor"
      role="img"
    >
      <title>{title}</title>
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M6 18L18 6M6 6l12 12"
      />
    </svg>
  );
}

function renderIcon(matched: boolean, t: ReturnType<typeof useTranslations>) {
  return matched ? (
    <CheckIcon title={t("complexity.matches")} />
  ) : (
    <CrossIcon title={t("complexity.doesNotMatch")} />
  );
}
const desc =
  "text-14px leading-4 text-input-light-label dark:text-input-dark-label";

export function PasswordComplexity({
  passwordComplexitySettings,
  password,
  equals,
}: Props) {
  const t = useTranslations("password");
  const hasMinLength = password?.length >= passwordComplexitySettings.minLength;
  const hasSymbol = symbolValidator(password);
  const hasNumber = numberValidator(password);
  const hasUppercase = upperCaseValidator(password);
  const hasLowercase = lowerCaseValidator(password);

  return (
    <div className="mb-4 grid grid-cols-2 gap-x-8 gap-y-2">
      {passwordComplexitySettings.minLength != undefined ? (
        <div className="flex flex-row items-center" data-testid="length-check">
          {renderIcon(hasMinLength, t)}
          <span className={desc}>
            <Translated
              i18nKey="complexity.length"
              namespace="password"
              data={{ minLength: passwordComplexitySettings.minLength.toString() }}
            />
          </span>
        </div>
      ) : (
        <span />
      )}
      <div className="flex flex-row items-center" data-testid="symbol-check">
        {renderIcon(hasSymbol, t)}
        <span className={desc}>
          <Translated i18nKey="complexity.hasSymbol" namespace="password" />
        </span>
      </div>
      <div className="flex flex-row items-center" data-testid="number-check">
        {renderIcon(hasNumber, t)}
        <span className={desc}>
          <Translated i18nKey="complexity.hasNumber" namespace="password" />
        </span>
      </div>
      <div className="flex flex-row items-center" data-testid="uppercase-check">
        {renderIcon(hasUppercase, t)}
        <span className={desc}>
          <Translated i18nKey="complexity.hasUppercase" namespace="password" />
        </span>
      </div>
      <div className="flex flex-row items-center" data-testid="lowercase-check">
        {renderIcon(hasLowercase, t)}
        <span className={desc}>
          <Translated i18nKey="complexity.hasLowercase" namespace="password" />
        </span>
      </div>
      <div className="flex flex-row items-center" data-testid="equal-check">
        {renderIcon(equals, t)}
        <span className={desc}>
          <Translated i18nKey="complexity.equals" namespace="password" />
        </span>
      </div>
    </div>
  );
}
