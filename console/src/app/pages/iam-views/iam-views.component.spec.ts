import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { IamViewsComponent } from './iam-views.component';

describe('IamViewsComponent', () => {
    let component: IamViewsComponent;
    let fixture: ComponentFixture<IamViewsComponent>;

    beforeEach(waitForAsync(() => {
        TestBed.configureTestingModule({
            declarations: [IamViewsComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(IamViewsComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
