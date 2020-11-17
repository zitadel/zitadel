import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CnslFormFieldComponent } from './form-field.component';

describe('CnslFormFieldComponent', () => {
    let component: CnslFormFieldComponent;
    let fixture: ComponentFixture<CnslFormFieldComponent>;

    beforeEach(async(() => {
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
