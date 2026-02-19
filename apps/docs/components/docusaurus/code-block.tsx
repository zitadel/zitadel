import { CodeBlock, Pre } from 'fumadocs-ui/components/codeblock';
import { highlight } from 'fumadocs-core/highlight';

export default async function CodeBlockWrapper({ children, language, title, className: _className, showLineNumbers: _showLineNumbers }: {
  children: string;
  language?: string;
  title?: string;
  className?: string;
  showLineNumbers?: boolean;
}) {
  const code = typeof children === 'string' ? children.trim() : '';
  const rendered = await highlight(code, {
    lang: language ?? 'text',
    themes: {
      light: 'github-light',
      dark: 'github-dark',
    },
    components: {
      pre: (props) => <Pre {...props} />,
    },
  });

  return (
    <CodeBlock title={title} keepBackground>
      {rendered}
    </CodeBlock>
  );
}
