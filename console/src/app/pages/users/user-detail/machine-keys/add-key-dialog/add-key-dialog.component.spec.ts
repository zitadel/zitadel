import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { AddKeyDialogComponent } from './add-key-dialog.component';

describe('AddKeyDialogComponent', () => {
    let component: AddKeyDialogComponent;
    let fixture: ComponentFixture<AddKeyDialogComponent>;

    beforeEach(waitForAsync(() => {
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
