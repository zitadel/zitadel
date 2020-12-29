import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { AuthUserDetailComponent } from './auth-user-detail.component';

describe('AuthUserDetailComponent', () => {
  let component: AuthUserDetailComponent;
  let fixture: ComponentFixture<AuthUserDetailComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [AuthUserDetailComponent],
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AuthUserDetailComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
