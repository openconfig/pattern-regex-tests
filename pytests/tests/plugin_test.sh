#!/bin/bash

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
REPO_DIR="$SCRIPT_DIR/../.."
cwd=$PWD
cd $REPO_DIR

pyang -p "testdata" --plugindir "pytests/plugins" --check-patterns "testdata/python-plugin-test.yang" 2>&1 | diff "pytests/tests/golden.txt" -
retcode=$?

cd $cwd
exit $retcode
