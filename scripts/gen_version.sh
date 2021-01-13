#!/bin/bash
latest=$(git describe --match "v[0-9]*[0-9]" --exclude "*rc*" --exclude "*beta*"  --exclude "*-*" --abbrev=0 HEAD | sed 's/v//')
beta=$(git describe --match "v[0-9]*beta" --exclude "*-*-*" --abbrev=0 HEAD | sed 's/v//')
rc=$(git describe --match "v[0-9]*rc[0-9]" --abbrev=0 | sed 's/v//')
echo -e "latest:$latest\nbeta:$beta\nrc:$rc" > version