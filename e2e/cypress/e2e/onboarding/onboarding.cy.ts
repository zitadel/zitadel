import { systemAuth } from "support/api/apiauth";
import { createInstance } from "support/api/instances";

it('should run on a virtual domain', () => {

    const domain = "mytestinstance13.e2e"
    systemAuth().then(system => {
        cy.request({
            url: `http://${domain}:8080/healthz`,
            failOnStatusCode: false
        }).then(res => {
            if (!res.isOkStatusCode) {
                createInstance(system, "testinstance", domain)
                cy.wait(10_000)
            }
        })
    })

    cy.context(`http://${domain}:8080`).as('ctx');
})