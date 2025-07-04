import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { PasswordDialogSMSProviderComponent } from './password-dialog-sms-provider.component';

describe('PasswordDialogComponent', () => {
  let component: PasswordDialogSMSProviderComponent;
  let fixture: ComponentFixture<PasswordDialogSMSProviderComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [PasswordDialogSMSProviderComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PasswordDialogSMSProviderComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
