#!/bin/sh

export CONF_FILE=${1}

touch $CONF_FILE

echo \
"flyway.url=jdbc:${DB_DRIVER}://${MIGRATION_DB_HOST}:${MIGRATION_DB_PORT}/${MIGRATION_DB_NAME}
flyway.user=$MIGRATION_DB_USER
flyway.password=$MIGRATION_DB_PASS" \
>> $CONF_FILE