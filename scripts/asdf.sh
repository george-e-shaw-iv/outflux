#!/usr/bin/env bash

# This script adds all of the plugins and installs of the tools that asdf manages for
# this repository.

if [[ ! -f .tool-versions ]]; then
    echo "This script needs to be ran in the root of the repository."
    exit 1
fi

# Add all of the plugins that are used.
asdf plugin add golang
asdf plugin add protoc
asdf plugin update protoc # This is necessary because older versions will not work on M1s.

# Install all of the tools listed in .tool-versions.
asdf install