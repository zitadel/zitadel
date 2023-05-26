const fs = require('fs')
const YAML = require('yaml')

const defaults = fs.readFileSync('../../cmd/defaults.yaml', 'utf8')
const steps = fs.readFileSync('../../cmd/setup/steps.yaml', 'utf8')


const docdata = [
  {
    key: "ZITADEL_FIRSTINSTANCE_ORG_HUMAN_USERNAME",
    desc: `In case that UserLoginMustBeDomain is false (default) and you don't overwrite the username with an email, it will be suffixed by the org domain (org-name + domain from config). for example: zitadel-admin in org ZITADEL on domain.tld -> zitadel-admin@zitadel.domain.tld`,
    type: undefined,
    shortlist: true,
  }
]


let doc = YAML.parseDocument(defaults)

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
    yield [...pre]
    // yield pre.join('_').toUpperCase()
}

let keys = Array.from(deepKeys(json))


const keys2env =  function(parts) {
  return "ZITADEL_" + parts.join("_").toUpperCase()
}

function getDocData(key, data) {
  let result = data.filter(obj => obj.key === key)
  if (result.length === 0) {
    return [null, null, null]
  }
  if (result.length > 1) {
    throw new Error("Duplicate Keys")
  }
  return [result[0]['desc'], result[0]['type'], result[0]['shortlist']]
}

let output = keys.map(parts => [
  keys2env(parts), 
  doc.getIn(parts), 
  ...getDocData(keys2env(parts), docdata),
])

console.log(output)

console.log(doc.contents.items[4].commentBefore)

console.log(doc.get('ExternalPort').toString())
console.log(doc.get('ExternalPort'))
console.log(doc.get('ExternalPort', true).comment)
console.log(doc.get('ExternalPort', true))

// Metrics
console.log(doc.contents.items)