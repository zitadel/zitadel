import { redirect } from "next/navigation";

export default function Page() {
  // automatically redirect to loginname
  if (!process.env.DEBUG) {
    redirect("/loginname");
  }
}
