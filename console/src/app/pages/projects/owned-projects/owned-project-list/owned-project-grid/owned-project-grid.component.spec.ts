import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { OwnedProjectGridComponent } from './owned-project-grid.component';

describe('GridComponent', () => {
    let component: OwnedProjectGridComponent;
    let fixture: ComponentFixture<OwnedProjectGridComponent>;

    beforeEach(waitForAsync(() => {
        TestBed.configureTestingModule({
            declarations: [OwnedProjectGridComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(OwnedProjectGridComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
