import parse_yaml from './src/parse_yaml.js'
import tablemark from 'tablemark'
import fs from 'fs'

const files = [
  {source: "../../../cmd/setup/steps.yaml", target: "_steps.mdx"},
  {source: "../../../cmd/defaults.yaml", target: "_defaults.mdx"}
]

files.map(file => {

  const fileContent = fs.readFileSync(file.source, 'utf8')

  const doc = parse_yaml(fileContent)

  // drop path variable and combine the comments
  const cleanDoc = doc.map(({path, ...item}) => [item.env, item.commentBefore + "\n" + item.comment, item.value])

  const markdown = tablemark(cleanDoc, {
      columns: [
        "Variable",
        "Comments",
        "If absent (default)"
      ],
      wrapWidth: 80,
    })

    try {
      fs.writeFileSync(file.target, markdown);
    } catch (err) {
      console.error(err);
    }

  console.log(markdown)
})