controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."

controller-gen "crd:trivialVersions=true" crd paths="./..." output:crd:artifacts:config=test
