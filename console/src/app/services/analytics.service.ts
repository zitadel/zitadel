import { Injectable } from "@angular/core";

@Injectable({
  providedIn: 'root',
})
export class AnalyticsService {

    public emitAnalyticsEvent(evt: Event, event_type: string, name: string, details: object ) {
        try {
            fetch('http://localhost:8080/events', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    event_data: { "event_type": event_type, "button_name": name, "details": details },
                    event: name
                }),
            }).catch(() => { });
        } catch { }
    }

}
