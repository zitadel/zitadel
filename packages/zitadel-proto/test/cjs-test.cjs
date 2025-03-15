// CommonJS import test
const zitadel = require("@zitadel/proto");

// Check if the import worked by accessing some properties
console.log("CommonJS import test:");
console.log("- Has v1 API:", !!zitadel.v1);
console.log("- Has v2 API:", !!zitadel.v2);
console.log("- Has v3alpha API:", !!zitadel.v3alpha);

// Test v1 API
console.log("- v1.user module:", !!zitadel.v1.user);
console.log("- v1.management module:", !!zitadel.v1.management);

// Test v2 API
console.log("- v2.user module:", !!zitadel.v2.user);
console.log("- v2.user_service module:", !!zitadel.v2.user_service);

// Test v3alpha API
console.log("- v3alpha.user module:", !!zitadel.v3alpha.user);
console.log("- v3alpha.user_service module:", !!zitadel.v3alpha.user_service);

// Test successful if we can access these modules
if (zitadel.v1 && zitadel.v2 && zitadel.v3alpha) {
  console.log("✅ CommonJS import test passed!");
} else {
  console.error("❌ CommonJS import test failed!");
  process.exit(1);
}
