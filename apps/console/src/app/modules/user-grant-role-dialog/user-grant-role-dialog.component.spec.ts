import { ComponentFixture, TestBed } from '@angular/core/testing';

import { UserGrantRoleDialogComponent } from './user-grant-role-dialog.component';

describe('UserGrantRoleDialogComponent', () => {
  let component: UserGrantRoleDialogComponent;
  let fixture: ComponentFixture<UserGrantRoleDialogComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [UserGrantRoleDialogComponent],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(UserGrantRoleDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
