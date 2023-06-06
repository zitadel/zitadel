import parse_yaml from './src/parse_yaml.js'
import tablemark from 'tablemark'
import fs from 'fs'

const files = [
  {source: "../../../cmd/setup/steps.yaml", target: "output/_steps.mdx"},
  {source: "../../../cmd/defaults.yaml", target: "output/_defaults.mdx"}
]

files.map(file => {

  const fileContent = fs.readFileSync(file.source, 'utf8')

  const doc = parse_yaml(fileContent)

 

  function combineComments(before, after) {
    let combined = before.trim() + (after !== '' ? '\n' + after.trim() : '')
    return combined.replace('\n ', '\n')
  }
  // console.log(doc.map(({...item}) => combineComments(item.commentBefore, item.comment)).filter(node => node !== ''))

  // drop path variable and combine the comments
  const cleanDoc = doc.map(({path, ...item}) => [item.env, combineComments(item.commentBefore, item.comment), String(item.value)])

  function toCellText (v) {
    return v
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/(?:\r\n|\r|\n)/g, '<br />');
  }

  const markdown = tablemark(cleanDoc, {
      columns: [
        "Variable",
        "Comments",
        "If absent (default)"
      ],
      // wrapWidth: 80,
      toCellText,
      lineEnding: "\n"
    })

    try {
      fs.writeFileSync(file.target, markdown);
    } catch (err) {
      console.error(err);
    }

  //console.log(markdown)
})