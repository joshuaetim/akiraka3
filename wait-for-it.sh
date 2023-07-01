#!/bin/sh
# wait-for-postgres.sh
  
# until mysql -u root -ppassword -e"quit"; do
#   >&2 echo "Postgres is unavailable - sleeping"
#   sleep 1
# done

sleep 10
  
# >&2 echo "Postgres is up - executing command"
# exec "$@"
./app/frontdesk
