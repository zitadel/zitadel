import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PasswordIamPolicyComponent } from './password-iam-policy.component';

describe('PasswordIamPolicyComponent', () => {
    let component: PasswordIamPolicyComponent;
    let fixture: ComponentFixture<PasswordIamPolicyComponent>;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            declarations: [PasswordIamPolicyComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(PasswordIamPolicyComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
