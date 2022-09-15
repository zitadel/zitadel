import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { ProjectPrivateLabelingDialogComponent } from './project-private-labeling-dialog.component';

describe('ProjectPrivateLabelingDialogComponent', () => {
  let component: ProjectPrivateLabelingDialogComponent;
  let fixture: ComponentFixture<ProjectPrivateLabelingDialogComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ProjectPrivateLabelingDialogComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ProjectPrivateLabelingDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
