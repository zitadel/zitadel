import React from "react";

export function SubProcessorTable() {

  const country_list = {
    us: "USA",
    eu: "EU",
  }
  const processors = [
    {
      entity: "Google LLC",
      services: "ZITADEL Cloud, Support Services",
      purpose: "Cloud infrastructure provider (Google Cloud), business applications and collaboration (Workspace), Data warehouse services, Content delivery network, DDoS and bot prevention",
      hosting: "Region designated by Customer, United States",
      country: country_list.us,
      enduserdata: true
    },
    {
      entity: "Cockroach Labs, Inc.",
      services: "ZITADEL Cloud",
      purpose: "Managed database services: Dedicated CockroachDB clusters on Google Cloud",
      hosting: "Region designated by Customer",
      country: country_list.us,
      enduserdata: true
    },
    {
      entity: "Datadog, Inc.",
      services: "ZITADEL Cloud, Websites",
      purpose: "",
      hosting: country_list.eu,
      country: country_list.us,
      enduserdata: true
    },
    {
      entity: "Github, Inc.",
      services: "ZITADEL",
      purpose: "Source code management, code scanning, dependency management, security advisory, issue management, continuous integration",
      hosting: country_list.us,
      country: country_list.us,
      enduserdata: false
    },
  ]

  return (
    <table className="text-xs">
      <tr>
        <th>Entity name</th>
        <th>Relevant services</th>
        <th>Purpose</th>
        <th>End-user data</th>
        <th>Hosting location</th>
        <th>Country of registration</th>
      </tr>
      {
        processors.map((processor, rowID) => {
          return (
            <tr>
              <td key={rowID}>{processor.entity}</td>
              <td>{processor.services}</td>
              <td>{processor.purpose}</td>
              <td>{processor.enduserdata ? 'Yes' : 'No'}</td>
              <td>{processor.hosting}</td>
              <td>{processor.country}</td>
            </tr>
          )
        })
      }
    </table>
  );
}
