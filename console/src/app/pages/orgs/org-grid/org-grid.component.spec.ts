import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { OrgGridComponent } from './org-grid.component';

describe('OrgGridComponent', () => {
    let component: OrgGridComponent;
    let fixture: ComponentFixture<OrgGridComponent>;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            declarations: [OrgGridComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(OrgGridComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
