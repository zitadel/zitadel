import React from "react";

export function SubProcessorTable() {

  const processors = [
    {
      entity: "Google",
      services: "ZITADEL Cloud",
      purpose: "We use Google as our infrastructure provider and for business applications and collaboration.",
      country: "United States"
    }
  ]
  let rows = "";

  for (const processor in processors) {
    rows = rows + `<tr><td>${processor.entity}</td><td>${processor.services}</td><td>${processor.purpose}</td><td>${processor.county}</td></tr>`
  }

  return (
    <table>
      <tr>
        <th>Entity name</th>
        <th>Relevant services</th>
        <th>Purpose</th>
        <th>Country of registration</th>
      </tr>
      {rows}
    </table>
  );
}
