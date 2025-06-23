"use client";

import { useTranslations } from "next-intl";
import { useRouter } from "next/navigation";
import { Button, ButtonVariants } from "./button";

export function BackButton() {
  const t = useTranslations("common");
  const router = useRouter();
  return (
    <Button
      onClick={() => router.back()}
      type="button"
      variant={ButtonVariants.Secondary}
    >
      {t("back")}
    </Button>
  );
}
