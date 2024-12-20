"use client";

import { Boundary } from "@/components/boundary";
import { Button } from "@/components/button";
import { LanguageProvider } from "@/components/language-provider";
import { useTranslations } from "next-intl";
import { useEffect } from "react";

export default function Error({ error, reset }: any) {
  useEffect(() => {
    console.log("logging error:", error);
  }, [error]);

  const t = useTranslations("error");

  return (
    <LanguageProvider>
      <Boundary labels={["Login Error"]} color="red">
        <div className="space-y-4">
          <div className="text-sm text-red-500 dark:text-red-500">
            <strong className="font-bold">Error:</strong> {error?.message}
          </div>
          <div>
            <Button onClick={() => reset()}>{t("tryagain")}</Button>
          </div>
        </div>
      </Boundary>
    </LanguageProvider>
  );
}
