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

  // If no items, it's likely an auto-generated list for the current sidebar category.
  // Fumadocs doesn't have a direct equivalent to "current sidebar category items" in a component easily accessible here without context.
  // For now, we'll render a placeholder or nothing to avoid build errors.
  // In a real migration, we'd use `useTreeContext` or similar if available, or rely on Fumadocs' index page generation.
  return (
    <div className="admonition admonition-info alert alert--info">
      <div className="admonition-heading">
        <h5>DocCardList Placeholder</h5>
      </div>
      <div className="admonition-content">
        <p>Auto-generated card lists are not yet fully supported in the migration.</p>
      </div>
    </div>
  );
}
