import parse_yaml from '../src/parse_yaml.js'

let file = 
`FirstInstance:
    # MachineKeyPath comment before
    MachineKeyPath:
    # Name of the first instance created
    InstanceName: ZITADEL # Default is ZITADEL
    SMTPConfiguration:
        # configuration of the host
        SMTP:
        # must include the port, like smtp.mailtrap.io:2525. IPv6 is also supported, like [2001:db8::1]:2525
        Host:
    Database:
        # CockroachDB is the default datbase of ZITADEL
        cockroach:
          Host: localhost
          Port: 26257
    MessageTexts:
        - MessageTextType: InitCode
          Language: de
          Title: Zitadel - User initialisieren
          PreHeader: User initialisieren
          Subject: User initialisieren
          Greeting: Hallo {{.DisplayName}},
          Text: Dieser Benutzer wurde soeben im Zitadel erstellt. Mit dem Benutzernamen &lt;br&gt;&lt;strong&gt;{{.PreferredLoginName}}&lt;/strong&gt;&lt;br&gt; kannst du dich anmelden. Nutze den untenstehenden Button, um die Initialisierung abzuschliessen &lt;br&gt;(Code &lt;strong&gt;{{.Code}}&lt;/strong&gt;).&lt;br&gt; Falls du dieses Mail nicht angefordert hast, kannst du es einfach ignorieren.
          ButtonText: Initialisierung abschliessen
        - MessageTextType: PasswordReset
          Language: de
          Title: Zitadel - Passwort zurücksetzen
          PreHeader: Passwort zurücksetzen
          Subject: Passwort zurücksetzen
          Greeting: Hallo {{.DisplayName}},
          Text: Wir haben eine Anfrage für das Zurücksetzen deines Passwortes bekommen. Du kannst den untenstehenden Button verwenden, um dein Passwort zurückzusetzen &lt;br&gt;(Code &lt;strong&gt;{{.Code}}&lt;/strong&gt;).&lt;br&gt; Falls du dieses Mail nicht angefordert hast, kannst du es ignorieren.
          ButtonText: Passwort zurücksetzen
    Quotas:
      # Items takes a slice of quota configurations, whereas for each unit type and instance, one or zero quotas may exist.
      # The following unit types are supported`


const doc = parse_yaml(file)

test('Expect two nodes', () => {
  expect(doc).toHaveLength(8);
});

test('Instance Name variable name', () => {
    expect(doc[1].env).toBe("ZITADEL_FIRSTINSTANCE_INSTANCENAME");
});

test('MachineKeyPath value should be null', () => {
    expect(doc[0].value).toBeNull();
});

test('Instance Name value', () => {
    expect(doc[1].value).toBe("ZITADEL");
});

test('Instance Name comment', () => {
    expect(doc[1].comment).toBe("Default is ZITADEL");
});

test('Instance Name description', () => {
    expect(doc[1].commentBefore).toBe("Name of the first instance created");
});

test('Comment before map', () => {
    expect(doc[0].commentBefore).toBe("MachineKeyPath comment before");
});

test('Array', () => {
    expect(doc[6].value).toBe("array[...]");
});

test('When no first value to map comment to, add as comment after (instead of before first item)', () => {
    expect(doc[7].comment).toBe(`Items takes a slice of quota configurations, whereas for each unit type and instance, one or zero quotas may exist.
 The following unit types are supported`);
});

test('Array', () => {
    expect(doc[6].value).toBe("array[...]");
});

console.log(doc)