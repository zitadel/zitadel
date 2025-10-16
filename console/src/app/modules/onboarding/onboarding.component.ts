import { Component } from '@angular/core';
import { AdminService } from 'src/app/services/admin.service';
import { ThemeService } from 'src/app/services/theme.service';
import { ONBOARDING_MILESTONES } from 'src/app/utils/onboarding';
import { AnalyticsService } from 'src/app/services/analytics.service';


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
    public analyticsService: AnalyticsService
  ) {
    console.log("OnboardingComponent constructor")
    this.adminService.loadMilestones.next(ONBOARDING_MILESTONES);
  }

}
