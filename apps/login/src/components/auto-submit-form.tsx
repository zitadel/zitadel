"use client";

import { useEffect, useRef } from "react";

type Props = {
  url: string;
  fields: Record<string, string>;
};

export function AutoSubmitForm({ url, fields }: Props) {
  const formRef = useRef<HTMLFormElement>(null);

  useEffect(() => {
    if (formRef.current) {
      formRef.current.submit();
    }
  }, []);

  return (
    <form ref={formRef} action={url} method="post">
      {Object.entries(fields).map(([key, value]) => (
        <input key={key} type="hidden" name={key} value={value} />
      ))}
      <noscript>
        <button type="submit">Continue</button>
      </noscript>
    </form>
  );
}
