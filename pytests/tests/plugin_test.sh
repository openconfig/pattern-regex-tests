#!/bin/bash

TEST_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
REPO_DIR="$TEST_DIR/../.."
cwd=$PWD
cd $TEST_DIR

pyang -p "$REPO_DIR/testdata" -p "testdata" --plugindir "$REPO_DIR/pytests/plugins" --check-patterns "testdata/python-plugin-test.yang" 2>&1 | diff - "golden.txt"
retcode=$?

cd $cwd
exit $retcode
