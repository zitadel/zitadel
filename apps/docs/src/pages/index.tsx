import { Button } from "@zitadel/core";
import { useIsomorphicLayoutEffect } from "@zitadel/utils";

export default function Docs() {
  useIsomorphicLayoutEffect(() => {
    console.log("zitadel docs page");
  }, []);
  return (
    <div>
      <h1>zitadel Documentation</h1>
      <Button>Click me</Button>
    </div>
  );
}
