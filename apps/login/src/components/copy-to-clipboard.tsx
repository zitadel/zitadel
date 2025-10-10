"use client";

import {
  ClipboardDocumentCheckIcon,
  ClipboardIcon,
} from "@heroicons/react/20/solid";
import copy from "copy-to-clipboard";
import { useEffect, useState } from "react";

type Props = {
  value: string;
};

export function CopyToClipboard({ value }: Props) {
  const [copied, setCopied] = useState(false);

  useEffect(() => {
    if (copied) {
      copy(value);
      const to = setTimeout(setCopied, 1000, false);
      return () => clearTimeout(to);
    }
  }, [copied]);

  return (
    <div className="flex flex-row items-center px-2">
      <button
        id="tooltip-ctc"
        type="button"
        className="text-primary-light-500 dark:text-primary-dark-500"
        onClick={() => setCopied(true)}
      >
        {!copied ? (
          <ClipboardIcon className="h-5 w-5" />
        ) : (
          <ClipboardDocumentCheckIcon className="h-5 w-5" />
        )}
      </button>
    </div>
  );
}
