import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { LabelPolicyComponent } from './label-policy.component';

describe('LabelPolicyComponent', () => {
    let component: LabelPolicyComponent;
    let fixture: ComponentFixture<LabelPolicyComponent>;

    beforeEach(waitForAsync(() => {
        TestBed.configureTestingModule({
            declarations: [LabelPolicyComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(LabelPolicyComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
