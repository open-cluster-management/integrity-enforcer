#!/bin/bash
#
# Copyright 2020 IBM Corporation.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -e

echo "IMAGE PUBLISH GOES HERE!"

echo "Pipeline for integrity verifier starting : $(date)"

# Tag images with COMPONENT_TAG_EXTENSION
#docker tag ${IV_SERVER_IMAGE_NAME_AND_VERSION} ${IV_SERVER_IMAGE_NAME_AND_VERSION}${COMPONENT_TAG_EXTENSION}
#docker tag ${IV_LOGGING_IMAGE_NAME_AND_VERSION} ${IV_LOGGING_IMAGE_NAME_AND_VERSION}${COMPONENT_TAG_EXTENSION}
#docker tag ${IV_OPERATOR_IMAGE_NAME_AND_VERSION} ${IV_OPERATOR_IMAGE_NAME_AND_VERSION}${COMPONENT_TAG_EXTENSION}

export COMPONENT_VERSION=${IV_VERSION}
export COMPONENT_DOCKER_REPO=${REGISTRY}

# Publish ${IV_IMAGE}
export COMPONENT_NAME=${IV_IMAGE}
export DOCKER_IMAGE_AND_TAG=${COMPONENT_DOCKER_REPO}/${COMPONENT_NAME}:${COMPONENT_VERSION}${COMPONENT_TAG_EXTENSION}
rm -rf pipeline
make pipeline-manifest/update PIPELINE_MANIFEST_COMPONENT_SHA256=${TRAVIS_COMMIT} PIPELINE_MANIFEST_COMPONENT_REPO=${TRAVIS_REPO_SLUG} PIPELINE_MANIFEST_BRANCH=${TRAVIS_BRANCH}
echo "Completed pipeline for integrity verifier component: $COMPONENT_NAME"

# Publish ${IV_LOGGING}
export COMPONENT_NAME=${IV_LOGGING}
export DOCKER_IMAGE_AND_TAG=${COMPONENT_DOCKER_REPO}/${COMPONENT_NAME}:${COMPONENT_VERSION}${COMPONENT_TAG_EXTENSION}
rm -rf pipeline
make pipeline-manifest/update PIPELINE_MANIFEST_COMPONENT_SHA256=${TRAVIS_COMMIT} PIPELINE_MANIFEST_COMPONENT_REPO=${TRAVIS_REPO_SLUG} PIPELINE_MANIFEST_BRANCH=${TRAVIS_BRANCH}
echo "Completed pipeline for integrity verifier component: $COMPONENT_NAME"

# Publish ${IV_OPERATOR}
export COMPONENT_NAME=${IV_OPERATOR}
export DOCKER_IMAGE_AND_TAG=${COMPONENT_DOCKER_REPO}/${COMPONENT_NAME}:${COMPONENT_VERSION}${COMPONENT_TAG_EXTENSION}
rm -rf pipeline
make pipeline-manifest/update PIPELINE_MANIFEST_COMPONENT_SHA256=${TRAVIS_COMMIT} PIPELINE_MANIFEST_COMPONENT_REPO=${TRAVIS_REPO_SLUG} PIPELINE_MANIFEST_BRANCH=${TRAVIS_BRANCH}
echo "Completed pipeline for integrity verifier component: $COMPONENT_NAME"

