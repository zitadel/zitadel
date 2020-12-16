import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { AuthPasswordlessComponent } from './auth-passwordless.component';

describe('AuthPasswordlessComponent', () => {
    let component: AuthPasswordlessComponent;
    let fixture: ComponentFixture<AuthPasswordlessComponent>;

    beforeEach(waitForAsync(() => {
        TestBed.configureTestingModule({
            declarations: [AuthPasswordlessComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(AuthPasswordlessComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
