#!/usr/bin/env bash
set -euo pipefail

# safe-data-report.sh — read-only data inventory across all ted1k sources
# Produces DATA-REPORT.md with current state of each data source
# All database access is SELECT-only via ephemeral docker containers
#
# TODO: rewrite as cmd/safe-data-report in Go, with URL-based connection config

REPORT="DATA-REPORT.md"

# --- queries (read-only, per-year range scan) ---

# $1=year — uses index range scan on stamp PK instead of full table scan
mysql_year_query() {
  cat <<SQL
SELECT $1, MIN(stamp), MAX(stamp), COUNT(*),
  TIMESTAMPDIFF(SECOND, MIN(stamp), MAX(stamp)) + 1 - COUNT(*),
  ROUND((TIMESTAMPDIFF(SECOND, MIN(stamp), MAX(stamp)) + 1 - COUNT(*))
    / (TIMESTAMPDIFF(SECOND, MIN(stamp), MAX(stamp)) + 1) * 100, 4)
FROM watt WHERE stamp >= '$1-01-01' AND stamp < '$(( $1 + 1 ))-01-01'
SQL
}

pg_year_query() {
  cat <<SQL
SELECT $1, MIN(stamp), MAX(stamp), COUNT(*),
  (EXTRACT(EPOCH FROM (MAX(stamp) - MIN(stamp))) + 1)::integer - COUNT(*),
  ROUND(((EXTRACT(EPOCH FROM (MAX(stamp) - MIN(stamp))) + 1)::integer - COUNT(*))::decimal
    / (EXTRACT(EPOCH FROM (MAX(stamp) - MIN(stamp))) + 1)::integer * 100, 2)
FROM watt WHERE stamp >= '$1-01-01' AND stamp < '$(( $1 + 1 ))-01-01'
SQL
}

# --- helpers ---

# Get the year range from a quick min/max query
mysql_run() {
  local host="$1" port="$2" user="$3" db="$4" query="$5"
  docker run --rm --network host --platform linux/amd64 mysql:5.7 \
    mysql -h "$host" -P "$port" -u "$user" -D "$db" \
    --batch --skip-column-names -e "$query" 2>/dev/null
}

pg_run() {
  local connstr="$1" query="$2"
  docker run --rm --network host postgres:14-alpine \
    psql "$connstr" -t -A -F $'\t' -c "$query" 2>/dev/null
}

summary_table_header() {
  echo "| Year | Min Date | Max Date | Count | Missing Samples | % Missing |"
  echo "| ---- | -------- | -------- | -----: | --------------: | --------: |"
}

summary_table_row() {
  # skip rows where count is 0 (empty year)
  [ "$4" = "0" ] && return
  printf "| %s | %s | %s | %s | %s | %s |\n" "$1" "$2" "$3" "$4" "$5" "$6"
}

report_mysql() {
  local label="$1" host="$2" port="$3" db="$4" user="$5" enabled="$6"
  echo ""
  echo "## ${label}"
  echo ""
  echo "Source: \`${user}@${host}:${port}/${db}\`"
  echo ""

  if [ "$enabled" != "true" ]; then
    echo "**Status**: SKIPPED (not yet active)"
    return
  fi

  if ! docker run --rm --network host --platform linux/amd64 mysql:5.7 \
    mysqladmin ping -h "$host" -P "$port" -u "$user" --silent 2>/dev/null; then
    echo "**Connectivity**: FAILED — host unreachable"
    return
  fi
  echo "**Connectivity**: OK"
  echo ""

  # Get year range from min/max (fast, uses index)
  local range
  range=$(mysql_run "$host" "$port" "$user" "$db" \
    "SELECT YEAR(MIN(stamp)), YEAR(MAX(stamp)) FROM watt") || true
  if [ -z "$range" ]; then
    echo "*Could not determine year range*"
    return
  fi

  local min_year max_year
  min_year=$(echo "$range" | cut -f1)
  max_year=$(echo "$range" | cut -f2)

  summary_table_header
  for year in $(seq "$min_year" "$max_year"); do
    local row
    row=$(mysql_run "$host" "$port" "$user" "$db" "$(mysql_year_query "$year")") || true
    if [ -n "$row" ]; then
      IFS=$'\t' read -r y mn mx cnt miss pct <<< "$row"
      summary_table_row "$y" "$mn" "$mx" "$cnt" "$miss" "$pct"
    fi
  done
}

report_pg() {
  local label="$1" host="$2" port="$3" db="$4" user="$5" pass="$6" enabled="$7"
  local connstr="postgresql://${user}:${pass}@${host}:${port}/${db}"
  echo ""
  echo "## ${label}"
  echo ""
  echo "Source: \`postgresql://${user}:***@${host}:${port}/${db}\`"
  echo ""

  if [ "$enabled" != "true" ]; then
    echo "**Status**: SKIPPED (not yet active)"
    return
  fi

  if ! docker run --rm --network host postgres:14-alpine \
    pg_isready -h "$host" -p "$port" -U "$user" 2>/dev/null; then
    echo "**Connectivity**: FAILED — host unreachable"
    return
  fi
  echo "**Connectivity**: OK"
  echo ""

  # Get year range from min/max (fast, uses index)
  local range
  range=$(pg_run "$connstr" \
    "SELECT EXTRACT(YEAR FROM MIN(stamp))::integer, EXTRACT(YEAR FROM MAX(stamp))::integer FROM watt") || true
  if [ -z "$range" ]; then
    echo "*Could not determine year range*"
    return
  fi

  local min_year max_year
  min_year=$(echo "$range" | cut -f1)
  max_year=$(echo "$range" | cut -f2)

  summary_table_header
  for year in $(seq "$min_year" "$max_year"); do
    local row
    row=$(pg_run "$connstr" "$(pg_year_query "$year")") || true
    if [ -n "$row" ]; then
      IFS=$'\t' read -r y mn mx cnt miss pct <<< "$row"
      summary_table_row "$y" "$mn" "$mx" "$cnt" "$miss" "$pct"
    fi
  done
}

report_jsonl_dir() {
  local label="$1" dir_path="$2"
  echo ""
  echo "## ${label}"
  echo ""
  echo "Path: \`${dir_path}/\`"
  echo ""

  if [ ! -d "$dir_path" ]; then
    echo "**Status**: Directory NOT FOUND"
    return
  fi

  local file_count
  file_count=$(find "$dir_path" -name '*.jsonl' -type f | wc -l | tr -d ' ')
  if [ "$file_count" -eq 0 ]; then
    echo "**Status**: Directory exists but no .jsonl files found"
    return
  fi

  local oldest newest total_size
  oldest=$(find "$dir_path" -name '*.jsonl' -type f | sort | head -1)
  newest=$(find "$dir_path" -name '*.jsonl' -type f | sort | tail -1)
  total_size=$(du -sh "$dir_path" | awk '{print $1}')

  echo "**Status**: ${file_count} files, ${total_size} total"
  echo ""
  echo "| | File |"
  echo "|---|---|"
  echo "| **Oldest** | $(basename "$oldest") |"
  echo "| **Newest** | $(basename "$newest") |"
}


# --- main ---

cat > "$REPORT" <<EOF
# Data Report

Generated by \`scripts/safe-data-report.sh\` — all queries are read-only.

- **Timestamp**: $(date -u +"%Y-%m-%dT%H:%M:%SZ")
- **Host**: $(hostname)
- **User**: $(whoami)
EOF

{
  # MySQL sources
  report_mysql "darwin/mysql (live source — capture data)" \
    darwin.imetrical.com 3306 ted root true

  report_mysql "dirac/mysql (future home — not yet active)" \
    dirac.imetrical.com 3306 ted root false

  # Postgres/TimescaleDB sources
  report_pg "d1-px1/timescaledb (live mirror)" \
    d1-px1.imetrical.com 5432 ted postgres secret true

  # JSONL
  report_jsonl_dir "jsonl — local data" \
    "./data/jsonl/month"

} >> "$REPORT"

echo "Report written to ${REPORT}"
