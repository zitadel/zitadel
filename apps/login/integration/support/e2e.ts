const url = Cypress.env("CORE_MOCK_STUBS_URL") || "http://localhost:22220/v1/stubs";

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
