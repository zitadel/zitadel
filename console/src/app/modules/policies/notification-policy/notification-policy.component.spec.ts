import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { NotificationPolicyComponent } from './notification-policy.component';

describe('PasswordComplexityPolicyComponent', () => {
  let component: NotificationPolicyComponent;
  let fixture: ComponentFixture<NotificationPolicyComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [NotificationPolicyComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NotificationPolicyComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
