-- hermes_fetch.sql — fetch Hermes work from ~/.hermes/state.db for the daily summary.
-- Schema: sessions(id, source[=channel], title, started_at/ended_at REAL epoch,
--   message_count, tool_call_count, input/output_tokens, estimated_cost_usd, ...)
--         messages(session_id, role, content, timestamp REAL epoch, tool_name, ...)
--
-- Time window: every query filters on the last 24h via strftime('%s','now','-24 hours').
-- To change the window, replace '-24 hours' (e.g. '-1 day', '-7 days', '-2 hours').
-- A message can land in the window even if its session started earlier, so activity
-- is detected by messages.timestamp, then attributed to a channel via the join.
--
-- Run a single query:  sqlite3 -box ~/.hermes/state.db < this_file.sql   (runs all)
-- Or copy one block:   sqlite3 -json ~/.hermes/state.db "<block>"

-- ============================================================================
-- Q1. Channel rollup in window — one row per source/channel.
-- ============================================================================
SELECT s.source                                   AS channel,
       count(DISTINCT m.session_id)               AS sessions,
       count(*)                                   AS messages,
       sum(m.role = 'user')                        AS user_msgs,
       sum(m.role = 'assistant')                   AS assistant_msgs,
       sum(m.tool_name IS NOT NULL AND m.tool_name <> '') AS tool_calls
FROM messages m
JOIN sessions s ON s.id = m.session_id
WHERE m.timestamp > strftime('%s','now','-24 hours')
GROUP BY s.source
ORDER BY messages DESC;

-- ============================================================================
-- Q2. Active sessions in window — title + channel + size + cost.
-- ============================================================================
SELECT s.source                                          AS channel,
       s.title,
       s.message_count                                   AS msgs,
       s.tool_call_count                                 AS tools,
       round(coalesce(s.actual_cost_usd, s.estimated_cost_usd, 0), 4) AS cost_usd,
       datetime(s.started_at, 'unixepoch', 'localtime')  AS started,
       s.end_reason
FROM sessions s
WHERE s.id IN (SELECT DISTINCT session_id FROM messages
               WHERE timestamp > strftime('%s','now','-24 hours'))
ORDER BY s.started_at DESC;

-- ============================================================================
-- Q3. Message bodies in window — user asks + assistant outcomes, by channel.
--     (content truncated to 300 chars; drop substr() for full text.)
-- ============================================================================
SELECT s.source                                          AS channel,
       datetime(m.timestamp, 'unixepoch', 'localtime')   AS ts,
       m.role,
       substr(m.content, 1, 300)                         AS content
FROM messages m
JOIN sessions s ON s.id = m.session_id
WHERE m.timestamp > strftime('%s','now','-24 hours')
  AND m.role IN ('user','assistant')
  AND length(coalesce(m.content,'')) > 0
ORDER BY m.timestamp;

-- ============================================================================
-- Q4. Message-channel sessions only (exclude local cron/cli/tui/subagent) —
--     the inbound work from slack/telegram/whatsapp/discord/signal/email/sms/...
-- ============================================================================
SELECT s.source                                          AS channel,
       s.user_id,
       s.title,
       s.message_count                                   AS msgs,
       datetime(s.started_at, 'unixepoch', 'localtime')  AS started
FROM sessions s
WHERE s.id IN (SELECT DISTINCT session_id FROM messages
               WHERE timestamp > strftime('%s','now','-24 hours'))
  AND s.source NOT IN ('cron','cli','tui','subagent','webui')
ORDER BY s.started_at DESC;

-- ============================================================================
-- Q5. Full-text search across messages (FTS5) — e.g. find a topic in any channel.
--     Replace 'voice' with your term; remove the time filter for all-time.
-- ============================================================================
SELECT s.source                                          AS channel,
       datetime(m.timestamp, 'unixepoch', 'localtime')   AS ts,
       m.role,
       snippet(messages_fts, 0, '[', ']', '…', 12)       AS hit
FROM messages_fts
JOIN messages m ON m.rowid = messages_fts.rowid
JOIN sessions s ON s.id = m.session_id
WHERE messages_fts MATCH 'voice'
  AND m.timestamp > strftime('%s','now','-24 hours')
ORDER BY m.timestamp DESC
LIMIT 20;
