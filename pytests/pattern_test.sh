#!/bin/bash

# Override OCDIR to test a different repository.
if [ -z $OCDIR ]; then
  OCDIR=$GOPATH/src/github.com/openconfig/public
  echo "\$OCDIR not given, using default: $OCDIR"
fi

TEST_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
REPO_DIR="$TEST_DIR/.."

pyang -p $OCDIR -p "$REPO_DIR/testdata" --plugindir "$REPO_DIR/pytests/plugins" --check-patterns "$REPO_DIR/testdata/regexp-test.yang"
