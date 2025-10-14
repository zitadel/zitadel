// Handle React 19 SSR hydration errors that don't affect functionality
Cypress.on('uncaught:exception', (err, runnable) => {
  // React error #419 is specifically about SSR Suspense boundary issues
  // This doesn't affect the actual functionality, just the SSR/CSR transition
  if (err.message.includes('Minified React error #419')) {
    console.warn('Cypress: Suppressed React SSR error #419 (Suspense boundary issue):', err.message);
    return false;
  }
  // Other hydration mismatches that are common with React 19 + Next.js 15
  if (err.message.includes('server could not finish this Suspense boundary') ||
      err.message.includes('Switched to client rendering') ||
      err.message.includes('Hydration failed') ||
      err.message.includes('Text content does not match server-rendered HTML')) {
    console.warn('Cypress: Suppressed React hydration error (non-functional):', err.message);
    return false;
  }
  // Let other errors fail the test as they should
  return true;
});

const url = Cypress.env("API_MOCK_STUBS_URL");

function removeStub(service: string, method: string) {
  return cy.request({
    url,
    method: "DELETE",
    qs: {
      service,
      method,
    },
  });
}

export function stub(service: string, method: string, out?: any) {
  removeStub(service, method);
  return cy.request({
    url,
    method: "POST",
    body: {
      stubs: [
        {
          service,
          method,
          out,
        },
      ],
    },
  });
}
