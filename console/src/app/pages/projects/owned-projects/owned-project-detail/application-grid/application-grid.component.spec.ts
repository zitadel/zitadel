import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { ApplicationGridComponent } from './application-grid.component';

describe('AppGridComponent', () => {
    let component: ApplicationGridComponent;
    let fixture: ComponentFixture<ApplicationGridComponent>;

    beforeEach(waitForAsync(() => {
        TestBed.configureTestingModule({
            declarations: [ApplicationGridComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(ApplicationGridComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
