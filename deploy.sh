#!/bin/bash

BRANCH="dev"

echo -e "\nSTEP 1: checking out to target branch and pulling the latest changes."

if [ $# -eq 1 ]; then
    BRANCH="$1"
else
    echo "INFO: No branch specified. Using the default branch '$BRANCH'."
fi

if ! git clean -xdf; then
  echo "ERROR: Failed to clean the untracked files present in a git working directory."
  exit 1
fi

if ! git reset --hard; then
  echo "ERROR: Failed to reset the Git branch. Aborting script."
  exit 1
fi

if ! git checkout "$BRANCH"; then
  echo "ERROR: Failed to checkout the $BRANCH. Aborting script."
  exit 1
fi

if ! git pull; then
  echo "ERROR: Failed to pull the latest changes from Git. Aborting script."
  exit 1
fi

echo -e "\nSTEP 3: Stop and remove any existing container with the same image."
EXISTING_CONTAINER=$(docker ps -q -f ancestor=hlsapi)
if [ "$EXISTING_CONTAINER" ]; then
    if ! docker stop "$EXISTING_CONTAINER"; then
      echo "ERROR: Failed to stop existing container, that shares the same image - '$EXISTING_CONTAINER'."
      exit 1
    fi
    if ! docker rm "$EXISTING_CONTAINER"; then
      echo "ERROR: Failed to remove existing container - '$EXISTING_CONTAINER'."
      exit 1
    fi
fi

EXISTING_IMAGE=$(docker images -q hlsapi)
if [ "$EXISTING_IMAGE" ]; then
    if ! docker rmi "$EXISTING_IMAGE"; then
      echo "ERROR: Failed to remove an old image - '$EXISTING_IMAGE'."
      exit 1
    fi
fi

echo -e "\nSTEP 6: Build the Docker image 'wasm-chat'."
if ! docker build -t hlsapi .; then
  echo "ERROR: Failed to build the Docker image."
  exit 1
fi

echo -e "\nSTEP 7: Run the Docker container with the new image and restart on failure."
if ! docker run -d --restart=always --network etha-chat --name hlsapi -p 9001:9001 hlsapi; then
  echo "ERROR: Failed to Run the Docker container."
  exit 1
fi

echo -e "\nSTEP 8: Making a deploy.sh executable."
if ! chmod +x deploy.sh; then
    echo "ERROR: Failed to make deploy.sh executable."
    exit 1
fi

echo "Removing unused Docker resources."
docker system prune -f