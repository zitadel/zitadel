import { dag, Container, Directory, object, func } from "@dagger.io/dagger"

@object()
// eslint-disable-next-line @typescript-eslint/no-unused-vars
class Console {

  @func()
  build(directory: Directory): Directory {
    return this.buildEnv(directory)
      .withDirectory("/src/", this.generate(directory))
      .withExec(["yarn", "run", "build"])
      .directory("./dist/console")
  }

  @func()
  generate(directory: Directory): Directory {
      return this.buildEnv(directory)
      .withExec(["yarn", "generate"])
      .withExec(["ls", "-la", "./src/app/proto/generated"])
      .directory("./src/app/proto/generated")
  }

  @func()
  buildEnv(directory: Directory): Container {
    return dag
    .container()
    .from("node:20")
    .withDirectory("/src/", directory, {include: ["console/**"]})
    .withDirectory("/src/", directory, {include: ["proto/**"]})
    .withDirectory("/src/", directory, {include: ["docs/frameworks.json"]})
    .withWorkdir("/src/console")
    .withMountedCache("/src/console/node_modules", dag.cacheVolume("console-node-modules"))
    .withExec(["yarn", "install", "--frozen-lockfile"])
  }

}
