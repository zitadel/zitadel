# login-standalone should be extended by the login-standalone target in apps/login/docker-bake.hcl
target "login-standalone" {
  dockerfile = "build/login/Dockerfile"
  cache-from = ["type=gha,scope=login-build-{{.Platform}}"]
  cache-to   = ["type=gha,mode=max,scope=login-build-{{.Platform}}"]
}

target "login-standalone-out" {
  inherits = ["login-standalone"]
  target   = "build-out"
  output   = ["type=local,dest=.artifacts/login"]
}
 