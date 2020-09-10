import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { AddIdpDialogComponent } from './add-idp-dialog.component';


describe('AddIdpDialogComponent', () => {
    let component: AddIdpDialogComponent;
    let fixture: ComponentFixture<AddIdpDialogComponent>;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            declarations: [AddIdpDialogComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(AddIdpDialogComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
