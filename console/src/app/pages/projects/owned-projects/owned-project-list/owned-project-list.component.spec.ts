import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { OwnedProjectListComponent } from './owned-project-list.component';

describe('OwnedProjectListComponent', () => {
    let component: OwnedProjectListComponent;
    let fixture: ComponentFixture<OwnedProjectListComponent>;

    beforeEach(waitForAsync(() => {
        TestBed.configureTestingModule({
            declarations: [OwnedProjectListComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(OwnedProjectListComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
