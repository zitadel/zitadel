import { Component } from '@angular/core';
import { AdminService, OnboardingActions } from 'src/app/services/admin.service';

export const ONBOARDING_EVENTS: OnboardingActions[] = [
  { order: 0, eventType: 'instance.policy.label.changed', link: ['/settings'], fragment: 'branding' },
  { order: 1, eventType: 'project.added', link: ['/projects/create'] },
  { order: 2, eventType: 'project.application.added', link: ['/projects/app-create'] },
  { order: 3, eventType: 'user.human.added', link: ['/users/create'] },
  { order: 4, eventType: 'instance.smtp.config.changed', link: ['/settings'], fragment: 'notifications' },
  { order: 5, eventType: 'user.grant.added', link: ['/grant-create'] },
  { order: 6, eventType: 'org.policy.notificaktion.added', link: '/settings?id=notifications' },
];

@Component({
  selector: 'cnsl-onboarding',
  templateUrl: './onboarding.component.html',
  styleUrls: ['./onboarding.component.scss'],
})
export class OnboardingComponent {
  public actions = this.adminService.progressEvents;

  constructor(public adminService: AdminService) {
    this.adminService.loadEvents.next(ONBOARDING_EVENTS.map((oe) => oe.eventType));
  }
}
