import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { ChangesComponent } from './changes.component';

describe('ChangesComponent', () => {
    let component: ChangesComponent;
    let fixture: ComponentFixture<ChangesComponent>;

    beforeEach(waitForAsync(() => {
        TestBed.configureTestingModule({
            declarations: [ChangesComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(ChangesComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
