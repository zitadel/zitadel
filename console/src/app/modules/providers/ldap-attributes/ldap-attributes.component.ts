import { Component, EventEmitter, Input, OnChanges, OnDestroy, Output } from '@angular/core';
import { FormControl, FormGroup } from '@angular/forms';
import { Subject, takeUntil } from 'rxjs';
import { LDAPAttributes } from 'src/app/proto/generated/zitadel/idp_pb';
import { requiredValidator } from '../../form-field/validators/validators';

@Component({
  selector: 'cnsl-ldap-attributes',
  templateUrl: './ldap-attributes.component.html',
  styleUrls: ['./ldap-attributes.component.scss'],
})
export class LDAPAttributesComponent implements OnChanges, OnDestroy {
  @Input() public initialAttributes?: LDAPAttributes.AsObject;
  @Output() public attributesChanged: EventEmitter<LDAPAttributes> = new EventEmitter<LDAPAttributes>();
  private destroy$: Subject<void> = new Subject();
  public form: FormGroup = new FormGroup({
    avatarUrlAttribute: new FormControl('', []),
    displayNameAttribute: new FormControl('', []),
    emailAttribute: new FormControl('', []),
    emailVerifiedAttribute: new FormControl('', []),
    firstNameAttribute: new FormControl('', []),
    idAttribute: new FormControl('', [requiredValidator]),
    lastNameAttribute: new FormControl('', []),
    nickNameAttribute: new FormControl('', []),
    phoneAttribute: new FormControl('', []),
    phoneVerifiedAttribute: new FormControl('', []),
    preferredLanguageAttribute: new FormControl('', []),
    preferredUsernameAttribute: new FormControl('', []),
    profileAttribute: new FormControl('', []),
  });

  constructor() {
    this.form.valueChanges.pipe(takeUntil(this.destroy$)).subscribe((value) => {
      if (value) {
        const attr = new LDAPAttributes();
        attr.setAvatarUrlAttribute(value.avatarUrlAttribute);
        attr.setDisplayNameAttribute(value.displayNameAttribute);
        attr.setEmailAttribute(value.emailAttribute);
        attr.setEmailVerifiedAttribute(value.emailVerifiedAttribute);
        attr.setFirstNameAttribute(value.firstNameAttribute);
        attr.setIdAttribute(value.idAttribute);
        attr.setLastNameAttribute(value.lastNameAttribute);
        attr.setNickNameAttribute(value.nickNameAttribute);
        attr.setPhoneAttribute(value.phoneAttribute);
        attr.setPhoneVerifiedAttribute(value.phoneVerifiedAttribute);
        attr.setPreferredLanguageAttribute(value.preferredLanguageAttribute);
        attr.setPreferredUsernameAttribute(value.preferredUsernameAttribute);
        attr.setProfileAttribute(value.profileAttribute);
        this.attributesChanged.emit(attr);
      }
    });
  }

  ngOnChanges(): void {
    if (this.initialAttributes) {
      this.form.patchValue(this.initialAttributes);
    }
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }
}
