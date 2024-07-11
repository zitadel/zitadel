"use client";

import Link from "next/link";
import { Button, ButtonVariants } from "./Button";

type Props = { hasBack?: boolean };

export default function BackButton({ hasBack }: Props) {
  return hasBack || history?.length > 1 ? (
    <Button
      onClick={() => history.back()}
      type="button"
      variant={ButtonVariants.Secondary}
    >
      back
    </Button>
  ) : null;
}
