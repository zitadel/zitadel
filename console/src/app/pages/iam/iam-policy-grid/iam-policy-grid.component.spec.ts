import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { IamPolicyGridComponent } from './iam-policy-grid.component';

describe('IamPolicyGridComponent', () => {
    let component: IamPolicyGridComponent;
    let fixture: ComponentFixture<IamPolicyGridComponent>;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            declarations: [IamPolicyGridComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(IamPolicyGridComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
