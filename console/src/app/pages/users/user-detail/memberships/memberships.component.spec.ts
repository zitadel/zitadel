import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { MembershipsComponent } from './memberships.component';

describe('MembershipsComponent', () => {
    let component: MembershipsComponent;
    let fixture: ComponentFixture<MembershipsComponent>;

    beforeEach(waitForAsync(() => {
        TestBed.configureTestingModule({
            declarations: [MembershipsComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(MembershipsComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
