import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { AppSecretDialogComponent } from './app-secret-dialog.component';

describe('AppSecretDialogComponent', () => {
    let component: AppSecretDialogComponent;
    let fixture: ComponentFixture<AppSecretDialogComponent>;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            declarations: [AppSecretDialogComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(AppSecretDialogComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
