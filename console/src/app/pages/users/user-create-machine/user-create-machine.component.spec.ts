import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { UserCreateMachineComponent } from './user-create-machine.component';

describe('UserCreateMachineComponent', () => {
    let component: UserCreateMachineComponent;
    let fixture: ComponentFixture<UserCreateMachineComponent>;

    beforeEach(waitForAsync(() => {
        TestBed.configureTestingModule({
            declarations: [UserCreateMachineComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(UserCreateMachineComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
