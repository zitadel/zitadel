import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { CnslFormFieldComponent } from './form-field.component';

describe('CnslFormFieldComponent', () => {
    let component: CnslFormFieldComponent;
    let fixture: ComponentFixture<CnslFormFieldComponent>;

    beforeEach(waitForAsync(() => {
        TestBed.configureTestingModule({
            declarations: [CnslFormFieldComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(CnslFormFieldComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
