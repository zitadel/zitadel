import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { PrivateLabelingPolicyComponent } from './private-labeling-policy.component';

describe('PrivateLabelingPolicyComponent', () => {
  let component: PrivateLabelingPolicyComponent;
  let fixture: ComponentFixture<PrivateLabelingPolicyComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [PrivateLabelingPolicyComponent],
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PrivateLabelingPolicyComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
