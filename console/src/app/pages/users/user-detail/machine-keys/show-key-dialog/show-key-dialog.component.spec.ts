import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { ShowKeyDialogComponent } from './show-key-dialog.component';

describe('ShowKeyDialogComponent', () => {
    let component: ShowKeyDialogComponent;
    let fixture: ComponentFixture<ShowKeyDialogComponent>;

    beforeEach(waitForAsync(() => {
        TestBed.configureTestingModule({
            declarations: [ShowKeyDialogComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(ShowKeyDialogComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
