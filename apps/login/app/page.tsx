import { demos } from "#/lib/demos";
import ThemeWrapper from "#/ui/ThemeWrapper";
import Link from "next/link";

export default function Page() {
  return (
    <div className="space-y-8">
      <h1 className="text-xl font-medium text-gray-800 dark:text-gray-300">
        Pages
      </h1>

      <div className="space-y-10 text-white">
        {demos.map((section) => {
          return (
            <div key={section.name} className="space-y-5">
              <div className="text-xs font-semibold uppercase tracking-wider text-gray-500">
                {section.name}
              </div>
              <div className="grid grid-cols-1 gap-5 lg:grid-cols-2">
                {section.items.map((item) => {
                  return (
                    <Link
                      href={`/${item.slug}`}
                      key={item.name}
                      className="bg-background-light-400 dark:bg-background-dark-400 group block space-y-1.5 rounded-lg px-5 py-3 hover:bg-background-light-500 hover:dark:bg-background-dark-300 hover:shadow-lg border border-gray-300 dark:border-gray-600 transition-all "
                    >
                      <div className="font-medium text-gray-600 dark:text-gray-200 group-hover:text-gray-900 dark:group-hover:text-gray-300">
                        {item.name}
                      </div>

                      {item.description ? (
                        <div className="line-clamp-3 text-sm text-gray-500 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-gray-300">
                          {item.description}
                        </div>
                      ) : null}
                    </Link>
                  );
                })}
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
