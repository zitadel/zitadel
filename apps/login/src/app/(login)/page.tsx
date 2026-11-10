import { redirect } from "next/navigation";

export default function Page() {
  // automatically redirect to loginname
  if (process.env.DEBUG !== "true") {
    redirect("/loginname");
  }
}
