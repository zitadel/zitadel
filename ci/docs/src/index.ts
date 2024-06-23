import { dag, Container, Directory, object, func } from "@dagger.io/dagger"

@object()
// eslint-disable-next-line @typescript-eslint/no-unused-vars
class Docs {

  @func()
  build(directory: Directory): Container {
    return dag
      .container()
      .from("node:20")
      .withWorkdir("/usr/local/docs")
      .withFile("/usr/local/docs/package.json", directory.file("docs/package.json"))
      .withFile("/usr/local/docs/yarn.lock", directory.file("docs/yarn.lock"))
      .withExec(["yarn", "install", "--frozen-lockfile"])
      .withExec(["npm", "cache", "clean", "--force"])
      .withExec(["mv", "/usr/local/docs/node_modules", "/node_modules"])
      .withDirectory("/usr/local/docs", directory, {include: ["docs"]})
      .withDirectory("/usr/local/proto", directory, {include: ["proto"]})
      .withExec(["yarn", "build"])
  }
}
