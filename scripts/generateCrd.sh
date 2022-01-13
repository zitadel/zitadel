controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./operator/..."

controller-gen crd paths="./operator/..." output:crd:artifacts:config=test
