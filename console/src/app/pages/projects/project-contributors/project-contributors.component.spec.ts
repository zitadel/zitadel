import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ProjectContributorsComponent } from './project-contributors.component';

describe('ProjectContributorsComponent', () => {
    let component: ProjectContributorsComponent;
    let fixture: ComponentFixture<ProjectContributorsComponent>;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            declarations: [ProjectContributorsComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(ProjectContributorsComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
