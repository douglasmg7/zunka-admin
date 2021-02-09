#!/usr/bin/env bash
MERCADO_LIVRE_USER_CODE=TG-60225cbe11a63d00064cbbe9-360790045

CMD="curl -X POST \
-H 'accept: application/json' \
-H 'content-type: application/x-www-form-urlencoded' \
'https://api.mercadolibre.com/oauth/token' \
-d 'grant_type=authorization_code' \
-d 'client_id=$MERCADO_LIVRE_APP_ID' \
-d 'client_secret=$MERCADO_LIVRE_SECRET_KEY' \
-d 'code=$MERCADO_LIVRE_USER_CODE' \
-d 'redirect_uri=$MERCADO_LIVRE_REDIRECT_URL'"

echo $CMD

eval $CMD
echo