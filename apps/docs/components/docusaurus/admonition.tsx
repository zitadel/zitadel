import { Callout } from 'fumadocs-ui/components/callout';

export default function Admonition({ type, title, children, icon }: any) {
  let calloutType: 'info' | 'warn' | 'error' = 'info';
  switch (type) {
    case 'note': calloutType = 'info'; break;
    case 'tip': calloutType = 'info'; break;
    case 'info': calloutType = 'info'; break;
    case 'warning': calloutType = 'warn'; break;
    case 'danger': calloutType = 'error'; break;
    case 'caution': calloutType = 'warn'; break;
  }
  
  return (
    <Callout title={title} type={calloutType} icon={icon}>
      {children}
    </Callout>
  );
}
