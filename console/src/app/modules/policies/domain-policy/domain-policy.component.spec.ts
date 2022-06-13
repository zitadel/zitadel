import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { DomainPolicyComponent } from './domain-policy.component';

describe('DomainPolicyComponent', () => {
  let component: DomainPolicyComponent;
  let fixture: ComponentFixture<DomainPolicyComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [DomainPolicyComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(DomainPolicyComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
