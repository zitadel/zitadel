import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PasswordComplexityPolicyComponent } from './password-complexity-policy.component';

describe('PasswordComplexityPolicyComponent', () => {
    let component: PasswordComplexityPolicyComponent;
    let fixture: ComponentFixture<PasswordComplexityPolicyComponent>;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            declarations: [PasswordComplexityPolicyComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(PasswordComplexityPolicyComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
