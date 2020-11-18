import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { ProjectGrantsComponent } from './project-grants.component';

describe('ProjectGrantsComponent', () => {
    let component: ProjectGrantsComponent;
    let fixture: ComponentFixture<ProjectGrantsComponent>;

    beforeEach(waitForAsync(() => {
        TestBed.configureTestingModule({
            declarations: [ProjectGrantsComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(ProjectGrantsComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
