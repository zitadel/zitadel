requestHttpRecord = {
  // all fields from baseStreamRecord
  "stream": "request_http",
  "version": "v1",
  "observed_time": "20250114T162059Z",
  "trace_id": "123",
  "span_id": "123",
  "module": "openid" | "saml" | "cache",
  "instance_id": "1234567890123",
  "org_id": "1234567890123",
  "user_id": "1234567890123",

  // all fields from requestBaseRecord
  "request_is_system_user": false,
  "request_is_authenticated": true,
  "request_duration": "50ms",

  "request_http_protocol": "",
  "request_http_host": "",
  "request_http_port": "",
  "request_http_path": "",
  "request_http_method": "",
  "request_http_status": 200,
  "request_http_referer": "",
  "request_http_user_agent": "",
  "request_http_remote_ip": "",
  "request_http_bytes_received": 1000,
  "request_http_bytes_sent": 1000,
}
