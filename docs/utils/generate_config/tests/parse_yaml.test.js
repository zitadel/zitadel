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
        User:
        Password:
        TLS:
        # if the host of the sender is different from ExternalDomain set DefaultInstance.DomainPolicy.SMTPSenderAddressMatchesInstanceDomain to false
        From:
        FromName:
    MessageTexts:
        - MessageTextType: InitCode`


const doc = parse_yaml(file)

test('Expect two nodes', () => {
  expect(doc).toHaveLength(2);
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

console.log(doc)