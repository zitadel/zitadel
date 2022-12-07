import * as React from "react";
import Prism from "prismjs";

export function Column({ title, language, children }) {
  const ref = React.useRef(null);

  React.useEffect(() => {
    if (ref.current) Prism.highlightElement(ref.current, false);
  }, [children]);

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div>{children}</div>
      <div className="overflow-hidden bg-white dark:bg-background-dark-400 border border-border-light dark:border-border-dark rounded-md w-full">
        <div className="py-2 px-4 bg-black/10 dark:bg-white/10">{title}</div>
        <div className="code" aria-live="polite">
          <pre ref={ref} className={`language-${language}`}>
            {`<div className="grid grid-cols-1 md:grid-cols-2 gap-4">
    <div>{children}</div>
    <div className="overflow-hidden bg-white dark:bg-background-dark-400 border border-border-light dark:border-border-dark rounded-md w-full">
    <div className="py-2 px-4 bg-black/10 dark:bg-white/10">
        {title}
    </div>
    <div>
        <code></code>
    </div>
    </div>
</div>
`}
          </pre>
        </div>
      </div>
      <style jsx>
        {`
          /* Override Prism styles */
          .code :global(pre[class*="language-"]) {
            margin: 0;
          }
        `}
      </style>
    </div>
  );
}
