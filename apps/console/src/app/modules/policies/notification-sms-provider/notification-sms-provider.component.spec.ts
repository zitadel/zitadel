import { ComponentFixture, TestBed } from '@angular/core/testing';

import { NotificationSMSProviderComponent } from './notification-sms-provider.component';

describe('NotificationSMSProviderComponent', () => {
  let component: NotificationSMSProviderComponent;
  let fixture: ComponentFixture<NotificationSMSProviderComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [NotificationSMSProviderComponent],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(NotificationSMSProviderComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
