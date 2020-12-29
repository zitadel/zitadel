import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { MachineKeysComponent } from './machine-keys.component';

describe('MachineKeysComponent', () => {
    let component: MachineKeysComponent;
    let fixture: ComponentFixture<MachineKeysComponent>;

    beforeEach(waitForAsync(() => {
        TestBed.configureTestingModule({
            declarations: [MachineKeysComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(MachineKeysComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
