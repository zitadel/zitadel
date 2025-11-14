import { useTranslations } from "next-intl";

export function Translated({
  i18nKey,
  namespace,
  data,
  ...props
}: {
  i18nKey: string;
  children?: React.ReactNode;
  namespace?: string;
  data?: any;
} & React.HTMLAttributes<HTMLSpanElement>) {
  const t = useTranslations(namespace);
  const helperKey = `${namespace ? `${namespace}.` : ""}${i18nKey}`;

  return (
    <span data-i18n-key={helperKey} {...props}>
      {t(i18nKey, data)}
    </span>
  );
}
