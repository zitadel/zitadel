import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { MemberCreateDialogComponent } from './member-create-dialog.component';


describe('AddMemberDialogComponent', () => {
    let component: MemberCreateDialogComponent;
    let fixture: ComponentFixture<MemberCreateDialogComponent>;

    beforeEach(async(() => {
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
