#!/bin/bash
latest=$(git describe --match "v[0-9]*" --exclude "*rc*" --exclude "*beta*"--abbrev=4 HEAD)
beta=$(git describe --match "v[0-9]*beta" --abbrev=4 HEAD)
rc=$(git describe --match "v[0-9]*rc*" --abbrev=4 HEAD)
echo -e "latest:$latest\nbeta:$beta\nrc:$rc" > version