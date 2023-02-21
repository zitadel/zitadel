import { Component, EventEmitter, Output } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import { AdminService } from 'src/app/services/admin.service';
import { StorageLocation, StorageService } from 'src/app/services/storage.service';

@Component({
  selector: 'cnsl-onboarding-card',
  templateUrl: './onboarding-card.component.html',
  styleUrls: ['./onboarding-card.component.scss'],
})
export class OnboardingCardComponent {
  public percentageChanged: EventEmitter<number> = new EventEmitter<number>();
  public loading$: BehaviorSubject<any> = new BehaviorSubject(false);
  public actions = this.adminService.progressEvents;
  @Output() public dismissedCard: EventEmitter<void> = new EventEmitter();

  constructor(public adminService: AdminService) {}

  public dismiss(): void {
    this.dismissedCard.emit();
  }
}
