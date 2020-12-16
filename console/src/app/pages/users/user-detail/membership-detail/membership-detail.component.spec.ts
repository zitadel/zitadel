import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { MembershipDetailComponent } from './membership-detail.component';

describe('MembershipDetailComponent', () => {
    let component: MembershipDetailComponent;
    let fixture: ComponentFixture<MembershipDetailComponent>;

    beforeEach(waitForAsync(() => {
        TestBed.configureTestingModule({
            declarations: [MembershipDetailComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(MembershipDetailComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
