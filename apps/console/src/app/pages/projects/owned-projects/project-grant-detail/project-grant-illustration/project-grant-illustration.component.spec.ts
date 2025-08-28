import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ProjectGrantIllustrationComponent } from './project-grant-illustration.component';

describe('ProjectGrantIllustrationComponent', () => {
  let component: ProjectGrantIllustrationComponent;
  let fixture: ComponentFixture<ProjectGrantIllustrationComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ProjectGrantIllustrationComponent],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ProjectGrantIllustrationComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
