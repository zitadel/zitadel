import { dag, Container, Directory, object, func } from "@dagger.io/dagger"

@object()
// eslint-disable-next-line @typescript-eslint/no-unused-vars
class Console {

  @func()
  build(directory: Directory): Container {
    return dag
      .container()
      .from("node:20")
      .withWorkdir("/usr/local/console")
      .withFile("/usr/local/console/package.json", directory.file("console/package.json"))
      .withFile("/usr/local/console/yarn.lock", directory.file("console/yarn.lock"))
      .withExec(["yarn", "install", "--frozen-lockfile"])
      .withExec(["npm", "cache", "clean", "--force"])
      .withExec(["mv", "/usr/local/console/node_modules", "/node_modules"])
      .withDirectory("/usr/local/console", directory, {include: ["console"]})
      .withDirectory("/usr/local/proto", directory, {include: ["proto"]})
      .withExec(["yarn", "generate"])
      .withExec(["yarn", "build"])
  }
}
