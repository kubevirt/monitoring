#!/usr/bin/env bash

set -ex

cp monitoringlinter testdata/src/a/testrepo
cd testdata/src/a/testrepo

[[ $(./monitoringlinter testrepo/... 2>&1 | wc -l) == 14 ]]
