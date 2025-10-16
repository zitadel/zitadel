import { Component, EventEmitter, OnInit, Output } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import { AdminService } from 'src/app/services/admin.service';
import { ONBOARDING_MILESTONES } from 'src/app/utils/onboarding';

@Component({
  selector: 'cnsl-onboarding-card',
  templateUrl: './onboarding-card.component.html',
  styleUrls: ['./onboarding-card.component.scss'],
  standalone: false,
})
export class OnboardingCardComponent implements OnInit {
  public percentageChanged: EventEmitter<number> = new EventEmitter<number>();
  public loading$: BehaviorSubject<any> = new BehaviorSubject(false);
  public actions = this.adminService.progressMilestones;
  @Output() public dismissedCard: EventEmitter<void> = new EventEmitter();

  constructor(public adminService: AdminService) {}

  public dismiss(): void {
    this.dismissedCard.emit();
  }

  public onRegisterClick(evt: Event, name: string, details: string| undefined) {
    // Fire-and-forget debug event; does not block navigation
    console.log("clicked onRegisterClick in OnboardingCardComponent")
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



  ngOnInit() {
    console.log("OnboardingCardComponent constructor ngOnInit")
    this.adminService.loadMilestones.next(ONBOARDING_MILESTONES);
  }
}
