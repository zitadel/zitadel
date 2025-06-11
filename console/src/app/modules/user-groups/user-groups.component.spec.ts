import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { UserGroupsComponent } from './user-groups.component';

describe('UserGroupsComponent', () => {
  let component: UserGroupsComponent;
  let fixture: ComponentFixture<UserGroupsComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [UserGroupsComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(UserGroupsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
