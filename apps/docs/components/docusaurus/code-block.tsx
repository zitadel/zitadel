import { codeToHtml } from 'shiki';

export default async function CodeBlock({ children, language, title, className: _className, showLineNumbers: _showLineNumbers }: {
  children: string;
  language?: string;
  title?: string;
  className?: string;
  showLineNumbers?: boolean;
}) {
  const code = typeof children === 'string' ? children.trim() : '';
  let highlighted: string;
  try {
    highlighted = await codeToHtml(code, {
      lang: language ?? 'text',
      themes: {
        light: 'github-light',
        dark: 'github-dark',
      },
      defaultColor: false,
    });
  } catch {
    highlighted = await codeToHtml(code, {
      lang: 'text',
      themes: {
        light: 'github-light',
        dark: 'github-dark',
      },
      defaultColor: false,
    });
  }

  return (
    <div className="my-4 overflow-hidden rounded-lg border border-fd-border">
      {title && (
        <div className="border-b border-fd-border bg-fd-muted px-4 py-2 text-sm font-medium text-fd-muted-foreground">
          {title}
        </div>
      )}
      <div
        className="[&_.shiki]:m-0 [&_.shiki]:overflow-x-auto [&_.shiki]:p-4 [&_.shiki]:text-sm [&_.shiki]:leading-relaxed [&_.shiki]:bg-(--shiki-light-bg) dark:[&_.shiki]:bg-(--shiki-dark-bg)"
        dangerouslySetInnerHTML={{ __html: highlighted }}
      />
    </div>
  );
}
