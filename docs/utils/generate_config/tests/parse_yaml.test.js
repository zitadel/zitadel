const parse_yaml = require('../src/parse_yaml');

let file = 
`FirstInstance:
    # MachineKeyPath comment before
    MachineKeyPath:
    # Name of the first instance created
    InstanceName: ZITADEL # Default is ZITADEL`


test('Expect two nodes', () => {
  expect(parse_yaml(file)).toHaveLength(2);
});

test('Instance Name variable name', () => {
    expect(parse_yaml(file)[1].env).toBe("ZITADEL_FIRSTINSTANCE_INSTANCENAME");
});

test('Instance Name value', () => {
    expect(parse_yaml(file)[1].value).toBe("ZITADEL");
});

test('Instance Name description', () => {
    expect(parse_yaml(file)[1].description).toBe("Name of the first instance created\nDefault is ZITADEL");
});

test('Comment before map', () => {
    expect(parse_yaml(file)[0].description).toBe("MachineKeyPath comment before");
});