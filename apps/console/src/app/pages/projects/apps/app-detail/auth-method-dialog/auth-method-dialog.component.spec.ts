import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { AuthMethodDialogComponent } from './auth-method-dialog.component';

describe('AuthMethodDialogComponent', () => {
  let component: AuthMethodDialogComponent;
  let fixture: ComponentFixture<AuthMethodDialogComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [AuthMethodDialogComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AuthMethodDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
