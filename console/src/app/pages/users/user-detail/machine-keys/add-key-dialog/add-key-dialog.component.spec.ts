import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { AddKeyDialogComponent } from './add-key-dialog.component';

describe('AddKeyDialogComponent', () => {
    let component: AddKeyDialogComponent;
    let fixture: ComponentFixture<AddKeyDialogComponent>;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            declarations: [AddKeyDialogComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(AddKeyDialogComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
