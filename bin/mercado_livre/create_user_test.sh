#!/usr/bin/env bash

CMD="curl -X POST \
-H 'Authorization: Bearer $MERCADO_LIVRE_ACCESS_TOKEN' \
-H 'Content-Type: application/json' \
'https://api.mercadolibre.com/users/test_user' \
-d '{ "site_id":"MLB" }'"

echo $CMD

eval $CMD
echo