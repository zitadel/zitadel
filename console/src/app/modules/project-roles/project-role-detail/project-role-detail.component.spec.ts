import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ProjectRoleDetailComponent } from './project-role-detail.component';

describe('ProjectRoleDetailComponent', () => {
    let component: ProjectRoleDetailComponent;
    let fixture: ComponentFixture<ProjectRoleDetailComponent>;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            declarations: [ProjectRoleDetailComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(ProjectRoleDetailComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
