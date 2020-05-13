import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PasswordPolicyComponent } from './password-policy.component';

describe('PasswordPolicyComponent', () => {
    let component: PasswordPolicyComponent;
    let fixture: ComponentFixture<PasswordPolicyComponent>;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            declarations: [PasswordPolicyComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(PasswordPolicyComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
