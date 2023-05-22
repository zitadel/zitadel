import { ExclamationTriangleIcon } from "@heroicons/react/24/outline";

type Props = {
  children: React.ReactNode;
};

export default function Alert({ children }: Props) {
  return (
    <div className="flex flex-row items-center justify-center border border-yellow-600/40 dark:border-yellow-500/20 bg-yellow-200/30 text-yellow-600 dark:bg-yellow-700/20 dark:text-yellow-200 rounded-md py-2 scroll-px-40">
      <ExclamationTriangleIcon className="h-5 w-5 mr-2" />
      <span className="text-center text-sm">{children}</span>
    </div>
  );
}
