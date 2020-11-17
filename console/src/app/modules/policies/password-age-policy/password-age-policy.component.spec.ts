import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { PasswordAgePolicyComponent } from './password-age-policy.component';

describe('PasswordAgePolicyComponent', () => {
    let component: PasswordAgePolicyComponent;
    let fixture: ComponentFixture<PasswordAgePolicyComponent>;

    beforeEach(waitForAsync(() => {
        TestBed.configureTestingModule({
            declarations: [PasswordAgePolicyComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(PasswordAgePolicyComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
