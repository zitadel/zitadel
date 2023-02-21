import { Component, EventEmitter } from '@angular/core';
import { BehaviorSubject, switchMap } from 'rxjs';
import { AuthServiceClient } from 'src/app/proto/generated/zitadel/AuthServiceClientPb';
import { Event } from 'src/app/proto/generated/zitadel/event_pb';
import { AdminService } from 'src/app/services/admin.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';

@Component({
  selector: 'cnsl-onboarding-card',
  templateUrl: './onboarding-card.component.html',
  styleUrls: ['./onboarding-card.component.scss'],
})
export class OnboardingCardComponent {
  public percentageChanged: EventEmitter<number> = new EventEmitter<number>();
  public loading$: BehaviorSubject<any> = new BehaviorSubject(false);
  public actions = this.adminService.progressEvents;
  public close: EventEmitter<void> = new EventEmitter();
  constructor(public adminService: AdminService) {}

  public dismiss(): void {
    this.close.emit();
  }
}
