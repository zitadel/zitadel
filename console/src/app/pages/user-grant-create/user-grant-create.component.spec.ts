import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { UserGrantCreateComponent } from './user-grant-create.component';

describe('UserGrantCreateComponent', () => {
    let component: UserGrantCreateComponent;
    let fixture: ComponentFixture<UserGrantCreateComponent>;

    beforeEach(waitForAsync(() => {
        TestBed.configureTestingModule({
            declarations: [UserGrantCreateComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(UserGrantCreateComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
