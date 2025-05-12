import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { GroupMemberCreateDialogComponent } from './group-member-create-dialog.component';

describe('AddMemberDialogComponent', () => {
  let component: GroupMemberCreateDialogComponent;
  let fixture: ComponentFixture<GroupMemberCreateDialogComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [GroupMemberCreateDialogComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(GroupMemberCreateDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
