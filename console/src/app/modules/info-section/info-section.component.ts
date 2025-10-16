import { ChangeDetectionStrategy, Component, Input } from '@angular/core';

export enum InfoSectionType {
  INFO = 'INFO',
  WARN = 'WARN',
  ALERT = 'ALERT',
}

@Component({
  selector: 'cnsl-info-section',
  templateUrl: './info-section.component.html',
  styleUrls: ['./info-section.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  standalone: false,
})
export class InfoSectionComponent {
  @Input() type: InfoSectionType = InfoSectionType.INFO;
  @Input() fitWidth: boolean = false;

  protected readonly infoSectionType = InfoSectionType;

  public onRegisterClick(evt: Event, name: string, details: string| undefined) {
    // Fire-and-forget debug event; does not block navigation
    console.log("clicked onRegisterClick in InfoSectionComponent")
    try {
      fetch('http://localhost:8080/events', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          event_data: {"event_type":"click", "button_name": name, details: details},
          instance_id: 'default', // TODO: pass real instance id if available in context
          parent_type: 'organization',
          parent_id: 'ORG_ID', // TODO: pass real org id if available
          table_name: 'projections.apps7',
          event: name,
        }),
      }).catch(() => {});
    } catch {}
  }


}
