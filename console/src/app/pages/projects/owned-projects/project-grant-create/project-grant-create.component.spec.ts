import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { ProjectGrantCreateComponent } from './project-grant-create.component';

describe('GrantCreateComponent', () => {
    let component: ProjectGrantCreateComponent;
    let fixture: ComponentFixture<ProjectGrantCreateComponent>;

    beforeEach(waitForAsync(() => {
        TestBed.configureTestingModule({
            declarations: [ProjectGrantCreateComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(ProjectGrantCreateComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
