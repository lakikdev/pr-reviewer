#!/usr/bin/env bash

# Exit immediately if a any command exits with a non-zero status
set -e
# Set destination branch
DEST_BRANCH=$1


# Create new pull request and get its ID
echo "Creating new branch: $DEST_BRANCH from $BITBUCKET_BRANCH"
if curl -X POST https://api.bitbucket.org/2.0/repositories/$BITBUCKET_REPO_OWNER/$BITBUCKET_REPO_SLUG/src \
  --fail --show-error --silent \
  --user $BITBUCKET_USERNAME:$BITBUCKET_PASSWORD \
  --form "parents=$BITBUCKET_BRANCH" --form "branch=$DEST_BRANCH" --form "message=Created by pipeline"; then
  printf '%s' "\nSuccess! Branch $DEST_BRANCH is created!"
else
  printf '%s' "\nERROR! Branch $DEST_BRANCH creating failed!"
  exit 1
fi;