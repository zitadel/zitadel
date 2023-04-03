import { Tab } from '#/ui/Tab';

export type Item = {
  text: string;
  slug?: string;
};

export const TabGroup = ({ path, items }: { path: string; items: Item[] }) => {
  return (
    <div className="-mt-2 flex flex-wrap items-center">
      {items.map((item) => (
        <Tab key={path + item.slug} item={item} path={path} />
      ))}
    </div>
  );
};
