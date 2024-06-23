import { dag, Container, Directory, object, func } from "@dagger.io/dagger"

@object()
// eslint-disable-next-line @typescript-eslint/no-unused-vars
class Console {

  @func()
  build(directory: Directory): Directory {
    return dag
      .container()
      .from("node:20")
      .withDirectory("/src/", directory, {include: ["console/**"]})
      .withDirectory("/src/", directory, {include: ["proto/**"]})
      .withDirectory("/src/", directory, {include: ["docs/frameworks.json"]})
      .withWorkdir("/src/console")
      //.withMountedCache("/src/console/node_modules", dag.cacheVolume("console-node-modules"))
      .withExec(["yarn", "install", "--frozen-lockfile"])
      .withExec(["yarn", "generate"])
      .withExec(["yarn", "run", "build"])
      .directory("./dist")
  }
}
