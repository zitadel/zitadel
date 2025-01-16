baseStreamRecord = {
  "stream": "request_http",
  "version": "v1",
  "observed_time": "20250114T162059Z", // ISO 8601 Timestamp
  // stream records that belong together can be joined using their trace_id
  // for example when debugging an API request, all request scoped stream records can be filtered by a trace_id.
  "trace_id": "123",
  // span_id is especially handy for debugging request latencies.
  "span_id": "123",
  "module": "openid" | "saml" | "cache", // The part of the program that produced the record
  "instance_id": "1234567890123", // If available
  "org_id": "1234567890123", // If available
  "user_id": "1234567890123", // If available
  // additional static properties if configured that should be present in all records.
  // for example
  // - "region": "US1",
  // - "runtime_service_version": "v2.67.2" // so it is not only available in a single runtime_service record but in all records
}

// The "runtime" stream contains normal log records that are written all over the Zitadel code.
// These fields are present in streams that have the runtime_ prefix
runtimeBaseStreamRecord = {
  ...baseStreamRecord, // all fields from baseStreamRecord

  "runtime_severity": "info",
  "runtime_message": "user created",
  // runtime_attributes_* contains additional information passed in the log record
  // these properties have no guaranteed schema.
  "runtime_attributes_userid": "1234567890123",
}

// if stream is runtime_service
// A record is only written once in a runtime lifecycle
runtimeServiceRecord = {
  ...runtimeBaseStreamRecord, // all fields from runtimeBaseStreamRecord

  "runtime_service_name": "zitadel",
  "runtime_service_version": "v2.67.2",
  "runtime_service_process": "sdsf321ew6f5", // For example Pod ID
}

// if stream is runtime_error
runtimeErrorRecord = {
  ...runtimeBaseStreamRecord, // all fields from runtimeBaseStreamRecord

  "runtime_error_cause": "user not found by email user@example.com: no rows in result set",
  "runtime_error_stack": "line1\nline2\nline3",
  "runtime_error_i18n_key": "Errors.User.NotFound", // If error is of type ZitadelError
  "runtime_error_type": "InternalError", // If error is of type ZitadelError
}

// These fields are present in streams that have the request_ prefix
requestBaseRecord = {
  ...baseStreamRecord, // all fields from baseStreamRecord

  "request_is_system_user": false,
  "request_is_authenticated": true,
  "request_duration": "50ms",
}

// if stream is request_http
requestHttpRecord = {
  ...requestBaseRecord, // all fields from requestBaseRecord

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

// if stream is request_grpc
requestGrpcRecord = {
  ...requestBaseRecord, // all fields from requestBaseRecord

  "request_grpc_service": "",
  "request_grpc_method": "",
  "request_grpc_code": "",
}

// These fields are present in streams that have the action_ prefix
actionTargetCallRecord = {
  "action_target_id": "",
  "action_name": "",
  "action_protocol": "",
  "action_host": "",
  "action_port": "",
  "action_path": "",
  "action_method": "",
  "action_status": 200,
}

// if stream is action_trigger_event
actionTriggerEventRecord = {
  "action_trigger_event_id": "",
}

// if stream is action_trigger_function
actionTriggerFunctionRecord = {
  "action_trigger_function_name": "",
}

// These fields are present in streams that have the action_trigger_grpc_ prefix
actionTriggerGrpcBaseRecord = {
  "action_trigger_grpc_service": "",
  "action_trigger_grpc_method": "",
}

// if stream is action_trigger_grpc_request
actionTriggerGrpcRequestRecord = {
  ...actionTriggerGrpcBaseRecord,  // all fields from actionTriggerGrpcBaseRecord
  // (no additional properties)
}

// if stream is action_trigger_grpc_response
actionTriggerGrpcResponseRecord = {
  ...actionTriggerGrpcBaseRecord,  // all fields from actionTriggerGrpcBaseRecord

  "action_trigger_grpc_response_code": 200,
}

// if stream is event
eventRecord = {
  "event_id": "",
  "event_sequence": "",
  "event_position": "",
  "event_type": "",
  "event_data": {}, // dynamically typed
  "event_editor_user": "",
  "event_version": "",
  "event_aggregate_id": "",
  "event_aggregate_type": "",
  "event_resource_owner": "",
}

// These fields are present in streams that have the notification_ prefix
notificationBaseRecord = {
  "notification_messagetype": "",
  "notification_triggering_event_id": "",
}

notificationEmailRecord = {
  ...notificationBaseRecord,   // all fields from notificationBaseRecord

  "notification_email_smtpprovider_id": "",
  "notification_email_receipient": "",
}

notificationSMSRecord = {
  ...notificationBaseRecord,   // all fields from notificationBaseRecord

  "notification_sms_phonenumber": "",
}

notificationWebhookRecord = {
  ...notificationBaseRecord,   // all fields from notificationBaseRecord

  "notification_webhook_url": "",
}
