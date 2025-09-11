import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { LoginTextsComponent } from './login-texts.component';

describe('LoginTextsComponent', () => {
  let component: LoginTextsComponent;
  let fixture: ComponentFixture<LoginTextsComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [LoginTextsComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(LoginTextsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
