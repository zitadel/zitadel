import { ComponentFixture, TestBed } from '@angular/core/testing';

import { LoginPolicyIdpsComponent } from './login-policy-idps.component';

describe('LoginPolicyIdpsComponent', () => {
  let component: LoginPolicyIdpsComponent;
  let fixture: ComponentFixture<LoginPolicyIdpsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [LoginPolicyIdpsComponent],
    })
      .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(LoginPolicyIdpsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
