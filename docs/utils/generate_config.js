const fs = require('fs')
const YAML = require('yaml')

const defaults = fs.readFileSync('../../cmd/defaults.yaml', 'utf8')
const steps = fs.readFileSync('../../cmd/setup/steps.yaml', 'utf8')

let doc = YAML.parseDocument(steps)

let json = doc.toJSON()
// console.log(JSON.stringify(json, null, 2))

//https://stackoverflow.com/a/47064979
function* deepKeys (t, pre = [])
{ if (Array.isArray(t))
    return
  else if (Object(t) === t)
    for (const [k, v] of Object.entries(t))
      yield* deepKeys(v, [...pre, k])
  else
    yield pre
    // yield pre.join('_').toUpperCase()
}

let keys = Array.from(deepKeys(json))



// let env_variables = keys.map(parts => keys2env(parts))

// console.log(env_variables)

// // ExternalPort example
// console.log(doc.get('ExternalPort').toString()) // String
// console.log(doc.get('ExternalPort')) // Number
// console.log(doc.get('ExternalPort', true).comment) // Comment on line
// console.log(doc.get('ExternalPort', true)) // Scalar

const keys2env =  function(parts, prefix = "ZITADEL_") {
  return prefix + parts.join("_").toUpperCase()
}

let output = []

let v = YAML.visit(doc, {
  Scalar(key, value, path) {

    // put the key names of the path in an array
    path_array = path.filter(node => YAML.isPair(node)).map(pair => pair.key.value)

    // we need to treat comments before a collection in such case
    // that it's attached to the first element instead
    let parent = path.slice(-2, -1)[0] // second to last element (aka. parent)
    
    if (key === 'key') {
      
      let description = ''
      // console.log(`FistItem: ${parent.items[0] === path.slice(-1)[0]} Value: ${path.slice(-1)[0]}`)
      if (parent.items[0] === path.slice(-1)[0] && parent.commentBefore !== undefined) description += parent.commentBefore // is a first item
      if(value.commentBefore !== undefined) description += value.commentBefore

      output.push({
        env: keys2env(path_array), 
        path: path_array,
        description: description.trim()
      })

    }

    if (key === 'value') {
      output_node = output.find(node => node.env === keys2env(path_array))
      output_node.value = value.value 
      if(value.comment !== undefined) output_node.description += value.comment
    }

    //if (parent.commentBefore !== undefined && parent.items[0]) output_node.commentBefore = parent.commentBefore

  }, 

})

  console.log(output.filter(nodes => nodes.value !== undefined))
// Metrics
// console.log(doc.contents.items[4].commentBefore)
//console.log(doc.contents.items)