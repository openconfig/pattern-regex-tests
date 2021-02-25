#!/bin/bash

# Override OCDIR to test a different repository.
if [ -z $OCDIR ]; then
  OCDIR=$GOPATH/src/github.com/openconfig/public
  echo "\$OCDIR not given, using default: $OCDIR"
fi

TEST_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
REPO_DIR="$TEST_DIR/.."

tmpstderr=$(mktemp)
pyang -p $OCDIR -p "$REPO_DIR/testdata" --msg-template="| {line} | {msg} |" --plugindir "$REPO_DIR/pytests/plugins" --check-patterns "$REPO_DIR/regexp-test/*.yang" 2> $tmpstderr
retcode=$?
if [ $retcode -ne 0 ]; then
  >&2 echo "| Line # | typedef | error |"
  >&2 echo "| --- | --- | --- |"
fi
>&2 cat $tmpstderr
rm $tmpstderr
exit $retcode
