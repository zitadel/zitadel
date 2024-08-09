import { redirect } from "next/navigation";

export default function Page() {
  // automatically redirect to loginname
  redirect("/loginname");
}
