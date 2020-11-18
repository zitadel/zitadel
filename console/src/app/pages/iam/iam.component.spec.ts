import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { IamComponent } from './iam.component';

describe('IamComponent', () => {
    let component: IamComponent;
    let fixture: ComponentFixture<IamComponent>;

    beforeEach(waitForAsync(() => {
        TestBed.configureTestingModule({
            declarations: [IamComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(IamComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
