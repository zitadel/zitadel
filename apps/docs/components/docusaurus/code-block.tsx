import React from 'react';

export default function CodeBlock({ children, language, title, className, showLineNumbers: _showLineNumbers }: any) {
  return (
    <div className="my-4">
      {title && <div className="bg-secondary text-secondary-foreground px-4 py-2 text-sm font-medium rounded-t-md">{title}</div>}
      <pre className={`p-4 overflow-x-auto bg-secondary/50 rounded-b-md ${title ? 'rounded-t-none' : 'rounded-md'} ${className || ''}`}>
        <code className={language ? `language-${language}` : ''}>
          {children}
        </code>
      </pre>
    </div>
  );
}
