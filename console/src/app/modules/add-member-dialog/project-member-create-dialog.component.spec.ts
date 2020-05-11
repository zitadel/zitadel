import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ProjectMemberCreateDialogComponent } from './project-member-create-dialog.component';


describe('AddMemberDialogComponent', () => {
    let component: ProjectMemberCreateDialogComponent;
    let fixture: ComponentFixture<ProjectMemberCreateDialogComponent>;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            declarations: [ProjectMemberCreateDialogComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(ProjectMemberCreateDialogComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
