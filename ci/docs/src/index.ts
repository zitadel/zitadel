import { dag, Container, Directory, object, func } from "@dagger.io/dagger"

@object()
// eslint-disable-next-line @typescript-eslint/no-unused-vars
class Docs {

  @func()
  build(directory: Directory): Container {
    return dag
      .container()
      .from("node:20")
      .withWorkdir("/usr/local/app")
      .withFile("/usr/local/app/package.json", directory.file("package.json"))
      .withFile("/usr/local/app/yarn.lock", directory.file("yarn.lock"))
      .withExec(["yarn", "install", "--frozen-lockfile"])
      .withExec(["npm", "cache", "clean", "--force"])
      .withExec(["mv", "/usr/local/app/node_modules", "/node_modules"])
      .withDirectory("/usr/local/app", directory)
  }
}
