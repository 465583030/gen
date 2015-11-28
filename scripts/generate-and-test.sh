#!/bin/bash

if [ -z "$1" ]; then
	exit
fi

export LLVM_VERSION=$1

# Generate the new Clang version
cd clang/ || exit

rm -rf clang-c/
rm *_gen.go

go-clang-gen || exit

cd ..

# Change versions in files
sed -i -e "s/3.4/${LLVM_VERSION}/g" .travis.yml
find . -type f -not -path '*/\.*' -exec sed -i -e "s/go-clang-phoenix-bootstrap/go-clang-phoenix-v${LLVM_VERSION}/g" {} +

# Install and test the version
make install || exit
make test || exit

# Show the current state of the repository
git status
