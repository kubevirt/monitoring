#!/usr/bin/env bash

set -ex

cp monitoringlocationlinter testdata/src/a/testrepo
cd testdata/src/a/testrepo

[[ $(./monitoringlocationlinter testrepo/... 2>&1 | wc -l) == 4 ]]
