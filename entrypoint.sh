#!/bin/sh

SWAGGER_ROOT=/opt/swagger

setup_swagger() {
  sed -i "s|https://petstore.swagger.io/v2/swagger.json|swagger.json|g" $SWAGGER_ROOT/swagger-initializer.js
  sed -i "s|http://example.com/api|$SWAGGER_ROOT/swagger.json|g" $SWAGGER_ROOT/swagger-initializer.js
}

setup_swagger

/data-manager migrate
/data-manager run