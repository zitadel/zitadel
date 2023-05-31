import parse_yaml from './src/parse_yaml.js'
import tablemark from 'tablemark'

let file = 
`FirstInstance:
    # MachineKeyPath comment before
    MachineKeyPath:
    # Name of the first instance created
    InstanceName: ZITADEL # Default is ZITADEL`


const doc = parse_yaml(file)

// drop path variable and combine the comments
const cleanDoc = doc.map(({path, ...item}) => [item.env, item.commentBefore + "\n" + item.comment, item.value])

const markdown = tablemark(cleanDoc, {
    columns: [
      "Variable Name",
      "Description",
      "Default Value"
    ]
  })

console.log(markdown)