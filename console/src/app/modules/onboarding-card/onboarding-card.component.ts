import { Component, EventEmitter } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import { AdminService } from 'src/app/services/admin.service';

@Component({
  selector: 'cnsl-onboarding-card',
  templateUrl: './onboarding-card.component.html',
  styleUrls: ['./onboarding-card.component.scss'],
})
export class OnboardingCardComponent {
  public percentageChanged: EventEmitter<number> = new EventEmitter<number>();
  public loading$: BehaviorSubject<any> = new BehaviorSubject(false);

  constructor(public adminService: AdminService) {}
}
