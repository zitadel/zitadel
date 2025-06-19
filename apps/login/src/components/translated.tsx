import { useTranslations } from "next-intl";

export function Translated({
  i18nKey,
  children,
  namespace,
  ...props
}: {
  i18nKey: string;
  children?: React.ReactNode;
  namespace?: string;
} & React.HTMLAttributes<HTMLSpanElement>) {
  const t = useTranslations(namespace);
  const helperKey = `${namespace ? `${namespace}.` : ""}${i18nKey}`;

  return (
    <span data-i18n-key={helperKey} {...props}>
      {t(i18nKey)}
    </span>
  );
}
