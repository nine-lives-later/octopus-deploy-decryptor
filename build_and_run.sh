#!/bin/bash

set -e

# check for go version manager
if [ -r "$HOME/.gvm/scripts/gvm" ]; then
  source "$HOME/.gvm/scripts/gvm"
fi

go build
go fmt ./... > /dev/null

# print the variables
pushd ./Octopus-Export-Example > /dev/null

../octopus-deploy-decryptor -password 'pass!w0rd'

popd > /dev/null
