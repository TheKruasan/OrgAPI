#!/bin/sh

DB_URL="postgres://postgres:postgres@postgres:5432/org_db?sslmode=disable"

case "$1" in
    status)
        goose -dir ./migrations postgres "$DB_URL" status
        ;;

    up)
        goose -dir ./migrations postgres "$DB_URL" up
        ;;

    down)
        goose -dir ./migrations postgres "$DB_URL" down
        ;;

    up-by-one)
        goose -dir ./migrations postgres "$DB_URL" up-by-one
        ;;

    *)
        echo "Usage:"
        echo "./migration.sh status"
        echo "./migration.sh up"
        echo "./migration.sh down"
        echo "./migration.sh up-by-one"
        ;;
esac