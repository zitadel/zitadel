import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { OwnedProjectsComponent } from './owned-projects.component';

describe('OwnedProjectComponent', () => {
    let component: OwnedProjectsComponent;
    let fixture: ComponentFixture<OwnedProjectsComponent>;

    beforeEach(waitForAsync(() => {
        TestBed.configureTestingModule({
            declarations: [OwnedProjectsComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(OwnedProjectsComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
