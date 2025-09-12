import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ProjectRoleChipComponent } from './project-role-chip.component';

describe('ProjectRoleChipComponent', () => {
  let component: ProjectRoleChipComponent;
  let fixture: ComponentFixture<ProjectRoleChipComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ProjectRoleChipComponent],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ProjectRoleChipComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
