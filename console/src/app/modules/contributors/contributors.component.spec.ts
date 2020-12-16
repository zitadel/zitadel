import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { ContributorsComponent } from './contributors.component';

describe('ContributorsComponent', () => {
    let component: ContributorsComponent;
    let fixture: ComponentFixture<ContributorsComponent>;

    beforeEach(waitForAsync(() => {
        TestBed.configureTestingModule({
            declarations: [ContributorsComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(ContributorsComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
