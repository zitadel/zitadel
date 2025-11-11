import { redirect } from "next/navigation";

export const revalidate = 3600; // 1 hour - revalidate cached data

export default function Page() {
  // automatically redirect to loginname
  if (process.env.DEBUG !== "true") {
    redirect("/loginname");
  }
}
