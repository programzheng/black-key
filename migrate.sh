if [ -z "$1" ]; then
    goose -h
    exit
fi

godotenv -f goose.env goose -dir ./migrations $1