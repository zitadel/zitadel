export function removeStub(service: string, method: string) {
  return cy.request({
    url: "http://localhost:22220/v1/stubs",
    method: "DELETE",
    qs: {
      service: service,
      method: method,
    },
  });
}

export function addStub(service: string, method: string, out?: any) {
  return cy.request({
    url: "http://localhost:22220/v1/stubs",
    method: "POST",
    body: {
      stubs: [
        {
          service: service,
          method: method,
          out: out,
        },
      ],
    },
  });
}
