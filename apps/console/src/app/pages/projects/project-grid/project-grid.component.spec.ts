import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { ProjectGridComponent } from './project-grid.component';

describe('GridComponent', () => {
  let component: ProjectGridComponent;
  let fixture: ComponentFixture<ProjectGridComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ProjectGridComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ProjectGridComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
