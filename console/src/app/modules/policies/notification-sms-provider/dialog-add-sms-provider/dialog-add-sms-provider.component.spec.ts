import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { DialogAddSMSProviderComponent } from './dialog-add-sms-provider.component';

describe('PasswordDialogComponent', () => {
  let component: DialogAddSMSProviderComponent;
  let fixture: ComponentFixture<DialogAddSMSProviderComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [DialogAddSMSProviderComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(DialogAddSMSProviderComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
