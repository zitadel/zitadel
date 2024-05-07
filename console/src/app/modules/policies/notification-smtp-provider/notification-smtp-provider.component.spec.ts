import { ComponentFixture, TestBed } from '@angular/core/testing';

import { NotificationSMTPProviderComponent } from './notification-smtp-provider.component';

describe('IdpSettingsComponent', () => {
  let component: NotificationSMTPProviderComponent;
  let fixture: ComponentFixture<NotificationSMTPProviderComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [NotificationSMTPProviderComponent],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(NotificationSMTPProviderComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
