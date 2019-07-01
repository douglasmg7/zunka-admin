#!/usr/bin/env bash
echo "Creating tables on bluewhale.db ..."
sqlite3 bluewhale.db < create_tables.sql