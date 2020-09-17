import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PasswordLockoutPolicyComponent } from './password-lockout-policy.component';

describe('PasswordLockoutPolicyComponent', () => {
    let component: PasswordLockoutPolicyComponent;
    let fixture: ComponentFixture<PasswordLockoutPolicyComponent>;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            declarations: [PasswordLockoutPolicyComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(PasswordLockoutPolicyComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
