import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { SmtpTestDialogComponent } from './smtp-test-dialog.component';

describe('SmtpTestDialogComponent', () => {
  let component: SmtpTestDialogComponent;
  let fixture: ComponentFixture<SmtpTestDialogComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [SmtpTestDialogComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SmtpTestDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
