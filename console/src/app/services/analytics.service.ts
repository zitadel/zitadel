import { Injectable } from "@angular/core";

@Injectable({
  providedIn: 'root',
})
export class AnalyticsService {

    
    public emitAnalyticsEvent(evt: Event, event_type: string, name: string, details: string | object | undefined) {
        // Fire-and-forget debug event; does not block navigation
        console.log('Emitting analytics event: ', name, details);
        try {
            fetch('http://localhost:8080/events', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    event_data: { "event_type": event_type, "button_name": name, "details": details },
                    // instance_id: instance_id, // TODO: pass real instance id if available in context                    
                    event: name
                }),
            }).catch(() => { });
        } catch { }
    }

}
