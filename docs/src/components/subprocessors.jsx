import React from "react";

export function SubProcessorTable() {

  const country_list = {
    us: "USA",
    eu: "EU",
    ch: "Switzerland",
    fr: "France"
  }
  const processors = [
    {
      entity: "Google LLC",
      purpose: "Cloud infrastructure provider (Google Cloud), business applications and collaboration (Workspace), Data warehouse services, Content delivery network, DDoS and bot prevention",
      hosting: "Region designated by Customer, United States",
      country: country_list.us,
      enduserdata: "Yes (transit)"
    },
    {
      entity: "Cockroach Labs, Inc.",
      purpose: "Managed database services: Dedicated CockroachDB clusters on Google Cloud",
      hosting: "Region designated by Customer",
      country: country_list.us,
      enduserdata: "Yes (at rest)"
    },
    {
      entity: "Datadog, Inc.",
      purpose: "Infrastructure monitoring, log analytics, and alerting",
      hosting: country_list.eu,
      country: country_list.us,
      enduserdata: "Yes (logs)"
    },
    {
      entity: "Github, Inc.",
      purpose: "Source code management, code scanning, dependency management, security advisory, issue management, continuous integration",
      hosting: country_list.us,
      country: country_list.us,
      enduserdata: false
    },
    {
      entity: "Stripe Payments Europe, Ltd.",
      purpose: "Subscription management, payment process",
      hosting: country_list.us,
      country: country_list.us,
      enduserdata: false
    },
    {
      entity: "Bexio AG",
      purpose: "Customer management, payment process",
      hosting: country_list.ch,
      country: country_list.ch,
      enduserdata: false
    },
    {
      entity: "Mailjet SAS",
      purpose: "Marketing automation",
      hosting: country_list.eu,
      country: country_list.fr,
      enduserdata: false
    },
    {
      entity: "AC PM LLC (Postmark)",
      purpose: "Transactional mails, if no customer owned SMTP service is configured",
      hosting: country_list.us,
      country: country_list.us,
      enduserdata: "Yes (opt-out)"
    },
    {
      entity: "Vercel, Inc.",
      purpose: "Website hosting",
      hosting: country_list.us,
      country: country_list.us,
      enduserdata: false
    },
    {
      entity: "Agolia SAS",
      purpose: "",
      hosting: country_list.us,
      country: undefined, // Not clear for OSS plan, sent a request to algolia
      enduserdata: false
    },
    {
      entity: "Discord Netherlands BV",
      purpose: "",
      hosting: country_list.us,
      country: country_list.us,
      enduserdata: false
    },
    {
      entity: "Statuspal",
      purpose: "",
      hosting: country_list.us,
      country: country_list.us,
      enduserdata: false
    },
    {
      entity: "Plausible",
      purpose: "",
      hosting: country_list.us,
      country: country_list.us,
      enduserdata: false
    },
    {
      entity: "Twillio",
      purpose: "",
      hosting: country_list.us,
      country: country_list.us,
      enduserdata: false
    },
    {
      entity: "Mohlmann Solutions",
      purpose: "",
      hosting: country_list.us,
      country: country_list.us,
      enduserdata: false
    },
    {
      entity: "Remote Sri Lanka",
      purpose: "",
      hosting: country_list.us,
      country: country_list.us,
      enduserdata: false
    },
    {
      entity: "Clickhouse",
      purpose: "",
      hosting: country_list.us,
      country: country_list.us,
      enduserdata: false
    },
  ]

  return (
    <table className="text-xs">
      <tr>
        <th>Entity name</th>
        <th>Purpose</th>
        <th>End-user data</th>
        <th>Hosting location</th>
        <th>Country of registration</th>
      </tr>
      {
        processors
          .sort((a, b) => {
            if (a.entity < b.entity) return -1
            if (a.entity > b.entity) return 1
            else return 0
          })
          .map((processor, rowID) => {
          return (
            <tr>
              <td key={rowID}>{processor.entity}</td>
              <td>{processor.purpose}</td>
              <td>{processor.enduserdata ? processor.enduserdata  : 'No'}</td>
              <td>{processor.hosting}</td>
              <td>{processor.country}</td>
            </tr>
          )
        })
      }
    </table>
  );
}
