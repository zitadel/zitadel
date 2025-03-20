import { OrganizationSchema as OrgSchema1 } from "@zitadel/proto/zitadel/org/v2/org_pb";
// FYI Reparsing as ES module because module syntax was detected. This incurs a performance overhead.
import { OrganizationSchema as OrgSchema2 } from "@zitadel/proto/zitadel/org/v2/org_pb.js";

console.log("Legacy import test:");
console.log("- Generated zitadel/org/v2/org_pb import (discouraged):", !!OrgSchema1);
console.log("- Generated zitadel/org/v2/org_pb.js import (recommended):", !!OrgSchema2);

// Test successful if we can access these modules and they are the same type
if (OrgSchema1 && OrgSchema2 && OrgSchema1 === OrgSchema2) {
  console.log("✅ Legacy import test passed!");
} else {
  console.error("❌ Legacy import test failed!");
  process.exit(1);
}
