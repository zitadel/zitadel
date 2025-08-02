# login-standalone should be extended by the login-standalone target in apps/login/docker-bake.hcl
target "login-standalone" {
  dockerfile = "build/login/Dockerfile"
}

target "login-standalone-out" {
  inherits = ["login-standalone"]
  target   = "build-out"
  output   = ["type=local,dest=.artifacts/login"]
}
 