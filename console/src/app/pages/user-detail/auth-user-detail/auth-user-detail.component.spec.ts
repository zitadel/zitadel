import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { AuthUserDetailComponent } from './auth-user-detail.component';

describe('AuthUserDetailComponent', () => {
  let component: AuthUserDetailComponent;
  let fixture: ComponentFixture<AuthUserDetailComponent>;

  beforeEach(async(() => {
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
