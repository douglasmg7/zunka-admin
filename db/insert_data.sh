#!/usr/bin/env bash
echo "Inserting data..."
sqlite3 zunka.db < data.sql
