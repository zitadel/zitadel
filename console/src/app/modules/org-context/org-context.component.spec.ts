import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { OrgContextComponent } from './org-context.component';

describe('OrgContextComponent', () => {
  let component: OrgContextComponent;
  let fixture: ComponentFixture<OrgContextComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [OrgContextComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(OrgContextComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
