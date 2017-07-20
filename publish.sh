#!/bin/bash
ORG="d3sw"
REPO="floop"

# TODO -> These will need to come from the build server.
USERNAME="d3scm"
TOKEN="c465671ebf06f86ed6c9726ff34ce42feeb4a8e1"

# Only publish the binaries if the version is of the form 'v#.#.#'.
if [ "$APP_VERSION" != "$APP_VERSION_SHORT" ]; then echo "$APP_VERSION is not a release version.  Skipping publish."; exit 0; fi

# Get the release-id for this repository and version.
RELEASE_ID=`curl -s -u "${USERNAME}:${TOKEN}" https://api.github.com/repos/${ORG}/${REPO}/releases/tags/v${APP_VERSION} | jq --raw-output '.id'`

# Only publish if the current version exists in GitHub as a release.
if [ "$RELEASE_ID" == "null" ]; then echo "Release not found! Skipping..."; exit 0; fi

# Upload each of the binaries.
echo "Found release $APP_VERSION. Uploading assets..."
for FILE in dist/*
do
    FILENAME=$(basename $FILE)
    if curl -si -u "${USERNAME}:${TOKEN}" -X POST -H 'Content-Type: application/zip' --data-binary "@${FILE}" "https://uploads.github.com/repos/${ORG}/${REPO}/releases/${RELEASE_ID}/assets?Content-Type=application/zip&name=${FILENAME}" | grep -q "201 Created"; then
        echo "$FILE uploaded."
    else
        echo "Failed to upload $FILE."
    fi
done