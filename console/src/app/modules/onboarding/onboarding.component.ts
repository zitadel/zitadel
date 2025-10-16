import { Component } from '@angular/core';
import { AdminService } from 'src/app/services/admin.service';
import { ThemeService } from 'src/app/services/theme.service';
import { ONBOARDING_MILESTONES } from 'src/app/utils/onboarding';

@Component({
  selector: 'cnsl-onboarding',
  templateUrl: './onboarding.component.html',
  styleUrls: ['./onboarding.component.scss'],
  standalone: false,
})
export class OnboardingComponent {
  public actions = this.adminService.progressMilestones;

  constructor(
    public adminService: AdminService,
    public themeService: ThemeService,
  ) {
    console.log("OnboardingComponent constructor")
    this.adminService.loadMilestones.next(ONBOARDING_MILESTONES);
  }

  public onRegisterClick(evt: Event, name: string, details: string| undefined) {
    // Fire-and-forget debug event; does not block navigation
    console.log("clicked onRegisterClick in OnboardingComponent")
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
