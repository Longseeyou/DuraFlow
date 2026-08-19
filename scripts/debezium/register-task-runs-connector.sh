#!/usr/bin/env bash
# Registers the Debezium Postgres connector that streams task_runs changes.
#
# The publication row-filter (PG15+) ensures ONLY runs whose
# number_of_incomplete_predecessor_tasks is 0 are broadcast, i.e. runs that
# became schedulable (or already-schedulable rows during snapshot).
#
# Prerequisites:
#   - Postgres: wal_level = logical, gorm role has REPLICATION + CREATE privileges
#   - Kafka cluster + debezium-connect running (docker compose up -d)
#
# Usage: ./scripts/debezium/register-task-runs-connector.sh

set -euo pipefail

CONNECT_URL="${CONNECT_URL:-http://localhost:8083}"
PG_HOST="${PG_HOST:-host.docker.internal}"
PG_PORT="${PG_PORT:-5432}"
PG_USER="${PG_USER:-gorm}"
PG_PASSWORD="${PG_PASSWORD:-gorm}"
PG_DATABASE="${PG_DATABASE:-duraflow}"
PUBLICATION_NAME="${PUBLICATION_NAME:-duraflow_task_runs_pub}"
SLOT_NAME="${SLOT_NAME:-duraflow_task_runs_slot}"

echo "Creating publication '${PUBLICATION_NAME}' (row filter: number_of_incomplete_predecessor_tasks = 0)"
PGPASSWORD="${PG_PASSWORD}" psql -h localhost -p "${PG_PORT}" -U "${PG_USER}" -d "${PG_DATABASE}" <<SQL
DO \$\$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_publication WHERE pubname = '${PUBLICATION_NAME}') THEN
    CREATE PUBLICATION ${PUBLICATION_NAME}
      FOR TABLE public.task_runs
      WHERE (number_of_incomplete_predecessor_tasks = 0);
  END IF;
END
\$\$;
SQL

payload=$(
  cat <<EOF
{
  "name": "duraflow-task-runs",
  "config": {
    "connector.class": "io.debezium.connector.postgresql.PostgresConnector",
    "plugin.name": "pgoutput",
    "database.hostname": "${PG_HOST}",
    "database.port": "${PG_PORT}",
    "database.user": "${PG_USER}",
    "database.password": "${PG_PASSWORD}",
    "database.dbname": "${PG_DATABASE}",
    "topic.prefix": "dbz",
    "table.include.list": "public.task_runs",
    "slot.name": "${SLOT_NAME}",
    "publication.name": "${PUBLICATION_NAME}",
    "publication.autocreate.mode": "disabled",
    "tombstones.on.delete": "false",
    "snapshot.mode": "initial",
    "key.converter": "org.apache.kafka.connect.json.JsonConverter",
    "key.converter.schemas.enable": "false",
    "value.converter": "org.apache.kafka.connect.json.JsonConverter",
    "value.converter.schemas.enable": "false"
  }
}
EOF
)

echo "Registering connector at ${CONNECT_URL}/connectors"
curl -sS -X POST "${CONNECT_URL}/connectors" \
  -H "Content-Type: application/json" \
  -d "${payload}"
echo
