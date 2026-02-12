CREATE TABLE IF NOT EXISTS spans (
  trace_id String,
  span_id String,
  parent_id String,
  operation String,
  kind String,
  start_time DateTime64(3),
  end_time DateTime64(3),
  status String,
  sampled UInt8,
  attributes JSON,
  events JSON,
  links JSON
)
ENGINE = MergeTree
PARTITION BY toDate(start_time)
ORDER BY (trace_id, start_time)
SETTINGS index_granularity = 8192;
