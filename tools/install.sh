#! /bin/sh

cd $GOPATH/src/github.com/caos/zitadel/tools
for imp in `cat tools.go | grep "-" | sed -E "s/_ \"(.*.+)\"/\1/g"`; do
	echo "installing $imp"
	go install $imp
done
cd -
