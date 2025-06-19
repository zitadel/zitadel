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

  return (
    <span data-i18n-key={i18nKey} {...props}>
      {t(i18nKey)}
    </span>
  );
}
