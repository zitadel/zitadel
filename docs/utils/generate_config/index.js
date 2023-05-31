import parse_yaml from './src/parse_yaml.js'
import tablemark from "tablemark"

let file = 
`FirstInstance:
    # MachineKeyPath comment before
    MachineKeyPath:
    # Name of the first instance created
    InstanceName: ZITADEL # Default is ZITADEL`


const doc = parse_yaml(file)

const markdown = tablemark(doc)

console.log(markdown)