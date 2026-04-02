import { redirect } from "next/navigation"

/**
 * Root page redirects to the overview dashboard.
 * In self-hosted mode, this is the main entry point.
 * In cloud mode, this could redirect to instance selection (future).
 */
export default function RootPage() {
  redirect("/overview")
}
