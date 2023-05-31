import YAML from 'yaml'

const keys2env =  function(parts, prefix = "ZITADEL_") {
    return prefix + parts.join("_").toUpperCase()
}

function parseZitadelYaml(file) {
    let output = []
    const doc = YAML.parseDocument(file)
    // first loop creates an array with all the env variables and the paths
    // this is required since we need to shift around some comments
    // and it's easier if we already have all possible keys to assign
    YAML.visit(doc, {
        Scalar(key, value, path) {
      
          // put the key names of the path in an array
          let path_array = path.filter(node => YAML.isPair(node)).map(pair => pair.key.value)
          let env = keys2env(path_array)
          
          if (key === 'key') {
            output.push({
              env: env, 
              path: path_array,
              comment: '',
              commentBefore: ''
            })
          }
      
          if (key === 'value') output.find(node => node.env === env).value = value.value
      
        }, 
      }
    )

    // filter out the nodes without a value assigned (parents)
    let keys = output.filter(nodes => nodes.value !== undefined)

    // loop through the envs and add comments
    keys.forEach(variable => {
    
      let pair = doc.getIn(variable.path, true)
      let index = keys.findIndex(key => key.env === variable.env)

      if(pair.comment !== undefined && variable.value !== null) {
        keys[index].comment = pair.comment.trim()
      }

      // this is a case where the comment is treated as inline comment
      // since the value of the Pair is NULL
      // imo this is a bug in the parsing library

      if(pair.comment !== undefined && variable.value === null) {
        keys[index+1].commentBefore = pair.comment.trim()
      }

    })

    // In this loop we have to check if the first comment is attached to a Map/Collection
    // and then put it on the first Pair instead
    YAML.visit(doc, {
      Scalar(key, value, path) {
    
        // put the key names of the path in an array
        let path_array = path.filter(node => YAML.isPair(node)).map(pair => pair.key.value)
        let env = keys2env(path_array)

        // we need to treat comments before a collection in such case
        // that it's attached to the first element instead
        let parent = path.slice(-2, -1)[0] // second to last element (aka. parent)
          
        if (key === 'key' && parent.items[0] === path.slice(-1)[0] && parent.commentBefore !== undefined)keys.find(node => node.env === env).commentBefore = parent.commentBefore.trim()  
      }, 
    
    })

    return keys // only key nodes
}
export default parseZitadelYaml;