# login-standalone can extend the login-standalone target in apps/login/docker-bake.hcl
target "login-standalone" {
  context = .
  dockerfile = dockerfiles/login.Dockerfile
}

target "login-standalone-out" {
  inherits = ["login-standalone"]
  target   = "build-out"
  output = [
    "type=local,dest=.artifacts/login"
  ]
}
