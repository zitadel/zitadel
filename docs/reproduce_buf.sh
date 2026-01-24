#!/bin/bash
mkdir -p temp_repro
cp -r ../proto/vTest/* temp_repro/
template_path=$(readlink -f buf.gen.yaml)
output_path=$(readlink -f openapi/debug_output)
rm -rf $output_path
mkdir -p $output_path

cd temp_repro
echo "Running buf in $(pwd)"
npx @bufbuild/buf generate --template $template_path --output $output_path --verbose
ls -R $output_path
cd ..
rm -rf temp_repro
