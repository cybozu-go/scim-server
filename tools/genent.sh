#!/bin/bash

# Because we want to be able to use any arbitrary version of the
# github.com/cybozu-go/scim module when running this script,
# we _DYNAMICALLY_ modify our requirements when needed.

# We need to do this because we need both the actual Go module
# and the helper file located in tools/cmd/genresources/objects.yaml

set -e

GENENT_TEMP=.genent-tmp
mkdir $GENENT_TEMP
GENENT_TEMP=$(cd $GENENT_TEMP; pwd -P) # absolute path

function cleanup() {
	echo "Cleaning up..."
	if [[ -e "$GENENT_TEMP" ]]; then
		rm -rf $GENENT_TEMP
	fi
}

trap cleanup EXIT

# Setup the SCIM directory
if [[ -z "$SCIM_DIR" ]]; then
	SCIM_DIR="$GENENT_TEMP/scim"
	git clone https://github.com/cybozu-go/scim $SCIM_DIR
fi

SCIM_DIR=$(cd $SCIM_DIR; pwd -P)

pushd tools/cmd/genent

ABS=$(pwd -P)

# Here comes the absolutely bad bad part
# First, make a copy of go.mod so that we can revert it back
# when required
cp go.mod "$GENENT_TEMP/go.mod"

go mod edit -replace=github.com/cybozu-go/scim=$(realpath --relative-to="$ABS" $SCIM_DIR)
go mod tidy
go build -o "$GENENT_TEMP/genent" main.go
go mod edit -dropreplace=github.com/cybozu-go/scim
go mod tidy
popd

$GENENT_TEMP/genent -clone-dir=$SCIM_DIR

