import { Tabs as FumadocsTabs, Tab as FumadocsTab } from 'fumadocs-ui/components/tabs';
import React, { isValidElement } from 'react';

export function TabItem({ value, label: _label, children, default: _isDefault, className, attributes }: any) {
  return <FumadocsTab value={value} className={className} {...attributes}>{children}</FumadocsTab>;
}

export default function Tabs({ children, defaultValue, values, groupId, className }: any) {
  let items: string[] = [];

  const childrenArray = React.Children.toArray(children);

  if (values) {
    items = values.map((v: any) => v.label);
  } else {
    childrenArray.forEach((child) => {
      if (isValidElement(child)) {
        const props = child.props as any;
        const label = props.label || props.value;
        items.push(label);
      }
    });
  }

  const newChildren = childrenArray.map((child, index) => {
    if (isValidElement(child)) {
      const props = child.props as any;
      let label = props.label;
      if (!label && values) {
        const valObj = values.find((v: any) => v.value === props.value);
        if (valObj) label = valObj.label;
      }
      if (!label) label = props.value;

      return <FumadocsTab key={index} value={label} {...props}>{props.children}</FumadocsTab>;
    }
    return child;
  });

  let newDefaultValue = defaultValue;
  if (defaultValue && values) {
    const valObj = values.find((v: any) => v.value === defaultValue);
    if (valObj) newDefaultValue = valObj.label;
  } else if (defaultValue) {
    const child = childrenArray.find((c: any) => (c as any).props.value === defaultValue);
    if (child) {
      newDefaultValue = (child as any).props.label || (child as any).props.value;
    }
  }

  return (
    <FumadocsTabs items={items} defaultValue={newDefaultValue} groupId={groupId} className={className}>
      {newChildren}
    </FumadocsTabs>
  );
}


