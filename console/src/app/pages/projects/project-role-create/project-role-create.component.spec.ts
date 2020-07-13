import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ProjectRoleCreateComponent } from './project-role-create.component';

describe('ProjectRoleCreateComponent', () => {
  let component: ProjectRoleCreateComponent;
  let fixture: ComponentFixture<ProjectRoleCreateComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ProjectRoleCreateComponent],
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ProjectRoleCreateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
