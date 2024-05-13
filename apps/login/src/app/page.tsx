import { demos } from "@/lib/demos";
import Link from "next/link";

export default function Page() {
  return (
    <div className="rounded-lg bg-vc-border-gradient dark:bg-dark-vc-border-gradient p-px shadow-lg shadow-black/5 dark:shadow-black/20 mb-10">
      <div className="rounded-lg bg-background-light-400 dark:bg-background-dark-500 px-8 py-12">
        <div className="space-y-8">
          <h1 className="text-xl font-medium">Pages</h1>

          <div className="space-y-10">
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
                          className="bg-background-light-400 dark:bg-background-dark-400 group block space-y-1.5 rounded-md px-5 py-3 hover:shadow-lg hover:dark:bg-white/10 border border-divider-light dark:border-divider-dark transition-all "
                        >
                          <div className="font-medium">{item.name}</div>

                          {item.description ? (
                            <div className="line-clamp-3 text-sm text-text-light-secondary-500 dark:text-text-dark-secondary-500">
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

          <div className="flex flex-col">
            <div className="mb-5 text-xs font-semibold uppercase tracking-wider text-gray-500">
              Deploy your own on Vercel
            </div>
            <a href="https://vercel.com/new/clone?repository-url=https%3A%2F%2Fgithub.com%2Fzitadel%2Ftypescript&root-directory=apps/login&env=ZITADEL_API_URL,ZITADEL_SERVICE_USER_ID,ZITADEL_SERVICE_USER_TOKEN&envDescription=Setup%20a%20service%20account%20with%20IAM_OWNER%20membership%20on%20your%20instance%20and%20provide%20its%20id%20and%20personal%20access%20token.&project-name=zitadel-login&repository-name=zitadel-login">
              <img src="https://vercel.com/button" alt="Deploy with Vercel" />
            </a>
          </div>
        </div>
      </div>
    </div>
  );
}
