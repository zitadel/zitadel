import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { LoginPolicyComponent } from './login-policy.component';

describe('LoginPolicyComponent', () => {
    let component: LoginPolicyComponent;
    let fixture: ComponentFixture<LoginPolicyComponent>;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            declarations: [LoginPolicyComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(LoginPolicyComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
