import { Component } from '@angular/core';
import { AdminService } from 'src/app/services/admin.service';
import { ThemeService } from 'src/app/services/theme.service';
import { ONBOARDING_MILESTONES } from 'src/app/utils/onboarding';

@Component({
  selector: 'cnsl-onboarding',
  templateUrl: './onboarding.component.html',
  styleUrls: ['./onboarding.component.scss'],
})
export class OnboardingComponent {
  public actions = this.adminService.progressMilestones;

  constructor(
    public adminService: AdminService,
    public themeService: ThemeService,
  ) {
    this.adminService.loadMilestones.next(ONBOARDING_MILESTONES);
  }
}
