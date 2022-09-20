import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { AddMemberRolesDialogComponent } from './add-member-roles-dialog.component';

describe('AddMemberRolesDialogComponent', () => {
  let component: AddMemberRolesDialogComponent;
  let fixture: ComponentFixture<AddMemberRolesDialogComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [AddMemberRolesDialogComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AddMemberRolesDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
