let webhookEventSchema = {
  sentStatus: 0,
  payload: {
    unit: 0,
    id: '',
    callURL: '',
    periodStart: new Date(),
    threshold: 0,
    usage: 0,
  },
};

export type ZITADELWebhookEvent = typeof webhookEventSchema;
