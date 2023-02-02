#!/bin/sh
CERT_DIR=/opt/postgres
CERT_PATH=${CERT_DIR}/ca.crt
echo "initializing ca certificate ${CERT_PATH}"
echo "${POSTGRES_CA_CERT}" > ${CERT_PATH}
export KC_HOSTNAME=localhost
export KC_HOSTNAME_PORT=8080
export KC_SPI_HOSTNAME_DEFAULT_ADMIN=localhost
export KC_DB=postgres
export KC_DB_URL="jdbc:postgresql://${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}"
export JDBC_PARAMS="sslmode=${POSTGRES_SSL_MODE}&sslrootcert=${CERT_PATH}"
export KC_DB_USERNAME=${POSTGRES_USERNAME}
export KC_DB_PASSWORD=${POSTGRES_PASSWORD}
echo "keycloak database url is ${KC_DB_URL}"
echo "starting keycloak..."
COMMAND="/opt/keycloak/bin/kc.sh $@"
echo "> ${COMMAND}"
${COMMAND}
