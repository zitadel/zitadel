import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { ProjectRoleDetailDialogComponent } from './project-role-detail-dialog.component';

describe('ProjectRoleDetailDialogComponent', () => {
  let component: ProjectRoleDetailDialogComponent;
  let fixture: ComponentFixture<ProjectRoleDetailDialogComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ProjectRoleDetailDialogComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ProjectRoleDetailDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
