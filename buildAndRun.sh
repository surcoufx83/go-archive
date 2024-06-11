#!/bin/sh

# Define variables
IMAGE_NAME="go-archive"
CONTAINER_NAME="go-archive-dev"
HOST_PORT=8080
CONTAINER_PORT=8080

# Build the Docker image
echo "Building the Docker image..."
docker build -t $IMAGE_NAME .

# Check if the container is running and stop it
if [ "$(docker ps -q -f name=$CONTAINER_NAME)" ]; then
    echo "Stopping and removing the old container..."
    docker stop $CONTAINER_NAME
    docker rm $CONTAINER_NAME
fi

# Run the new container
echo "Starting the new container..."
docker run -d -p $HOST_PORT:$CONTAINER_PORT --name $CONTAINER_NAME $IMAGE_NAME

echo "Deployment complete. The container is running and accessible at http://localhost:$HOST_PORT"
