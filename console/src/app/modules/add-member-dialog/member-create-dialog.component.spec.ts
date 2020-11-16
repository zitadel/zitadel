import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { MemberCreateDialogComponent } from './member-create-dialog.component';


describe('AddMemberDialogComponent', () => {
    let component: MemberCreateDialogComponent;
    let fixture: ComponentFixture<MemberCreateDialogComponent>;

    beforeEach(waitForAsync(() => {
        TestBed.configureTestingModule({
            declarations: [MemberCreateDialogComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(MemberCreateDialogComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
