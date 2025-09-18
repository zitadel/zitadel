import { Suspense } from "react";
import { ApiReferenceComponent } from "@/components/ApiReference";

export default function HomePage() {
  return (
    <main>
      <Suspense fallback={<div>Loading API documentation...</div>}>
        <ApiReferenceComponent />
      </Suspense>
    </main>
  );
}
