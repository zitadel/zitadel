import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { PasswordlessComponent } from './passwordless.component';

describe('AuthPasswordlessComponent', () => {
  let component: PasswordlessComponent;
  let fixture: ComponentFixture<PasswordlessComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [PasswordlessComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PasswordlessComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
