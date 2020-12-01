#!/bin/bash

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
REPO_DIR="$SCRIPT_DIR/.."

pyang -p $OCDIR -p "$REPO_DIR/testdata" --plugindir "$REPO_DIR/pytests/plugins" --check-patterns "$REPO_DIR/testdata/regexp-test.yang"
