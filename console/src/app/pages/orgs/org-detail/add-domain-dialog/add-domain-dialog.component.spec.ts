import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { AddDomainDialogComponent } from './add-domain-dialog.component';

describe('WarnDialogComponent', () => {
    let component: AddDomainDialogComponent;
    let fixture: ComponentFixture<AddDomainDialogComponent>;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            declarations: [AddDomainDialogComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(AddDomainDialogComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
