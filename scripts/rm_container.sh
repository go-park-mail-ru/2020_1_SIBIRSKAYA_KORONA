CONTAINER_LIST=$(docker ps -a -q)

if [ -z "$CONTAINER_LIST" ]
then 
    echo "No containers"
else
    docker stop $CONTAINER_LIST
    docker rm $CONTAINER_LIST
fi