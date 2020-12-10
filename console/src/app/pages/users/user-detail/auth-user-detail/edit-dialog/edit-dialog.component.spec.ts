import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { EditDialogComponent } from './edit-dialog.component';

describe('CodeDialogComponent', () => {
    let component: EditDialogComponent;
    let fixture: ComponentFixture<EditDialogComponent>;

    beforeEach(waitForAsync(() => {
        TestBed.configureTestingModule({
            declarations: [EditDialogComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(EditDialogComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
