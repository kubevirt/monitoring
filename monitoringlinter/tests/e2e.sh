#!/usr/bin/env bash

set -ex

cp ./bin/monitoringlinter monitoringlinter/testdata/src/a/testrepo
cd monitoringlinter/testdata/src/a/testrepo

[[ $(./monitoringlinter testrepo/... 2>&1 | wc -l) == 14 ]]
