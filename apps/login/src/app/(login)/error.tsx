"use client";

import { Boundary } from "@/components/boundary";
import { Button } from "@/components/button";
import { Translated } from "@/components/translated";
import { useEffect } from "react";

export default function Error({ error, reset }: any) {
  useEffect(() => {
    console.log("logging error:", error);
  }, [error]);

  return (
    <Boundary labels={["Login Error"]} color="red">
      <div className="space-y-4">
        <div className="text-sm text-red-500 dark:text-red-500">
          <strong className="font-bold">Error:</strong> {error?.message}
        </div>
        <div>
          <Button data-i18n-key="error.tryagain" onClick={() => reset()}>
            <Translated i18nKey="tryagain" namespace="error" />
          </Button>
        </div>
      </div>
    </Boundary>
  );
}
