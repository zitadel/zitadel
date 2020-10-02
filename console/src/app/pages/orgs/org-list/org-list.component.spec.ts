import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { OrgListComponent } from './org-list.component';

describe('OrgListComponent', () => {
    let component: OrgListComponent;
    let fixture: ComponentFixture<OrgListComponent>;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            declarations: [OrgListComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(OrgListComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
