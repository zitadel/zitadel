import { Component, Input, OnInit } from '@angular/core';
import { FormControl } from '@angular/forms';
import { PasswordComplexityPolicy } from 'src/app/proto/generated/zitadel/management_pb';

@Component({
    selector: 'app-password-complexity-view',
    templateUrl: './password-complexity-view.component.html',
    styleUrls: ['./password-complexity-view.component.scss'],
})
export class PasswordComplexityViewComponent implements OnInit {
    @Input() public password!: FormControl;
    @Input() public policy!: PasswordComplexityPolicy.AsObject;
    constructor() { }

    ngOnInit(): void {
    }

}
