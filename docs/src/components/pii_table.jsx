import React from "react";

export function PiiTable() {

  const pii = [
    {
      type: "Basic data",
      examples: [
        'Names',
        'Email addresses',
        'User names'
      ],
      subjects: "All users as uploaded by Customer."
    },
    {
      type: "Login data",
      examples: [
        'Randomly generated ID',
        'Passwords',
        'Public keys / certificates ("FIDO2", "U2F", "x509", ...)',
        'User names or identifiers of external login providers',
        'Phone numbers',
      ],
      subjects: "All users as uploaded and feature use by Customer."
    },
    {
      type: "Profile data",
      examples: [
        'Profile pictures',
        'Gender',
        'Languages',
        'Nicknames or Display names',
        'Phone numbers',
        'Metadata'
      ],
      subjects: "All users as uploaded by Customer"
    },
    {
      type: "Communication data",
      examples: [
        'Emails',
        'Chats',
        'Call metadata',
        'Call recording and transcripts',
        'Form submissions',
      ],
      subjects: "Customers and users who communicate with us directly (e.g. support, chat)."
    },
    {
      type: "Payment data",
      examples: [
        'Billing address',
        'Payment information',
        'Customer number',
        'Support Customer history',
        'Credit rating information',
      ],
      subjects: "Customers who use services that require payment. Credit rating information: Only customers who pay by invoice."
    },
    {
      type: "Analytics data",
      examples: [
        'Usage metrics',
        'User behavior',
        'User journeys (eg, Milestones)',
        'Telemetry data',
        'Client-side anonymized session replay',
      ],
      subjects: "Customers who use our services."
    },
    {
      type: "Usage meta data",
      examples: [
        'User agent',
        'IP addresses',
        'Operating system',
        'Time and date',
        'URL',
        'Referrer URL',
        'Accepted Language',
      ],
      subjects: "All users"
    },
    ]

  return (
    <table className="text-xs">
      <tr>
        <th>Type of personal data</th>
        <th>Examples</th>
        <th>Affected data subjects</th>
      </tr>
      {
        pii.map((row, rowID) => {
          return (
            <tr>
              <td key={rowID}>{row.type}</td>
              <td><ul>{row.examples.map((example) => { return ( <li>{example}</li> )})}</ul></td>
              <td>{row.subjects}</td>
            </tr>
          )
        })
      }
    </table>
  );
}
