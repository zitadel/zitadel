"use client";

import { Button, ButtonVariants } from "./Button";

export default function BackButton() {
  return history && history.length > 1 ? (
    <Button
      onClick={() => history.back()}
      type="button"
      variant={ButtonVariants.Secondary}
    >
      back
    </Button>
  ) : null;
}
