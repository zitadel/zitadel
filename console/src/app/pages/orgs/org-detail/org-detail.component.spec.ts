import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { OrgDetailComponent } from './org-detail.component';

describe('OrgDetailComponent', () => {
  let component: OrgDetailComponent;
  let fixture: ComponentFixture<OrgDetailComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [OrgDetailComponent],
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(OrgDetailComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
