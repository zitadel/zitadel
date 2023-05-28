const YAML = require('yaml')

const keys2env =  function(parts, prefix = "ZITADEL_") {
    return prefix + parts.join("_").toUpperCase()
}

function parseZitadelYaml(file) {
    let output = []
    const doc = YAML.parseDocument(file)
    YAML.visit(doc, {
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
            // When the previous Scalar value is null (no value in the yaml), 
            // then the comment will be appended to this scalar. This seems
            // like a bug, but we can handle it by checking if the value is null
            if(value.comment !== undefined && value.value !== null) output_node.description += value.comment.trim()

            // now this comment needs to be appended to the next key !
            if(value.comment !== undefined && value.value === null) {
              next_key = parent.items[parent.items.findIndex(n => n === path.slice(-1)[0])+1].key.value
              path_array_next = [ ...path_array.slice(0, -1), next_key]
              console.log(path_array_next)
              console.log(keys2env(path_array_next))
            }
          }
      
          //if (parent.commentBefore !== undefined && parent.items[0]) output_node.commentBefore = parent.commentBefore
      
        }, 
      
      })

      return output.filter(nodes => nodes.value !== undefined) // only key nodes
}

module.exports = parseZitadelYaml;