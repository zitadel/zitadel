import React from 'react';
import { Cards, Card } from 'fumadocs-ui/components/card';

export default function DocCardList({ items }: { items?: any[] }) {
  // If items are provided (manual list), render them
  if (items) {
    return (
      <Cards>
        {items.map((item, index) => (
          <Card key={index} href={item.href} title={item.label} />
        ))}
      </Cards>
    );
  }

  return null;
}
