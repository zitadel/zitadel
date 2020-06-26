import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { OwnedProjectGridComponent } from './owned-project-grid.component';

describe('GridComponent', () => {
    let component: OwnedProjectGridComponent;
    let fixture: ComponentFixture<OwnedProjectGridComponent>;

    beforeEach(async(() => {
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
