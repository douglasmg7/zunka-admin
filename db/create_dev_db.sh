#!/usr/bin/env bash
echo "Creating tables on bluewhale.db ..."
sqlite3 bluewhale.db < create_tables.sql

echo "Populating tables on bluewhale.db ..."
sqlite3 bluewhale.db < populate_tables.sql
