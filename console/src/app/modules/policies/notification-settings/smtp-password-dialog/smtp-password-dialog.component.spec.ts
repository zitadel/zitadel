import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { SMTPPasswordDialogComponent } from './smtp-password-dialog.component';

describe('CodeDialogComponent', () => {
  let component: SMTPPasswordDialogComponent;
  let fixture: ComponentFixture<SMTPPasswordDialogComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [SMTPPasswordDialogComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SMTPPasswordDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
