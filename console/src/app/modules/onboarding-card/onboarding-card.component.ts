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

  ngOnInit() {
    this.adminService.loadMilestones.next(ONBOARDING_MILESTONES);
  }
}
