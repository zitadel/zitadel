import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { OwnedProjectContributorsComponent } from './owned-project-contributors.component';

describe('ProjectContributorsComponent', () => {
    let component: OwnedProjectContributorsComponent;
    let fixture: ComponentFixture<OwnedProjectContributorsComponent>;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            declarations: [OwnedProjectContributorsComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(OwnedProjectContributorsComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
