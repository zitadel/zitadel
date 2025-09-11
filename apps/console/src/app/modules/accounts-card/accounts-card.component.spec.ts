import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { AccountsCardComponent } from './accounts-card.component';

describe('AccountsCardComponent', () => {
  let component: AccountsCardComponent;
  let fixture: ComponentFixture<AccountsCardComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [AccountsCardComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AccountsCardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
