BAD_IMAGES_LIST=$(docker images -f "dangling=true" -q)

if [ -z "$BAD_IMAGES_LIST" ]
then 
    echo "No bad images"
else
    docker rmi $BAD_IMAGES_LIST
fi