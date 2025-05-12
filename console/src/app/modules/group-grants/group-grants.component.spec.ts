import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { GroupGrantsComponent } from './group-grants.component';

describe('GroupGrantsComponent', () => {
  let component: GroupGrantsComponent;
  let fixture: ComponentFixture<GroupGrantsComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [GroupGrantsComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(GroupGrantsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
