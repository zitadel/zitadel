import { Component } from '@angular/core';
import { AdminService } from 'src/app/services/admin.service';
import { ONBOARDING_EVENTS } from 'src/app/utils/onboarding';

@Component({
  selector: 'cnsl-onboarding',
  templateUrl: './onboarding.component.html',
  styleUrls: ['./onboarding.component.scss'],
})
export class OnboardingComponent {
  public actions = this.adminService.progressEvents;

  constructor(public adminService: AdminService) {
    this.adminService.loadEvents.next(ONBOARDING_EVENTS);
  }
}
