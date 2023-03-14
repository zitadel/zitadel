import { Component, EventEmitter, Input, OnChanges, OnDestroy, Output } from '@angular/core';
import { AbstractControl, UntypedFormBuilder, UntypedFormGroup } from '@angular/forms';
import { MatLegacyDialog as MatDialog } from '@angular/material/legacy-dialog';
import { Subscription } from 'rxjs';
import { requiredValidator } from 'src/app/modules/form-field/validators/validators';
import { Gender, Human, Profile } from 'src/app/proto/generated/zitadel/user_pb';

import { ProfilePictureComponent } from './profile-picture/profile-picture.component';

@Component({
  selector: 'cnsl-detail-form',
  templateUrl: './detail-form.component.html',
  styleUrls: ['./detail-form.component.scss'],
})
export class DetailFormComponent implements OnDestroy, OnChanges {
  @Input() public showEditImage: boolean = false;
  @Input() public preferredLoginName: string = '';
  @Input() public username!: string;
  @Input() public user!: Human.AsObject;
  @Input() public disabled: boolean = true;
  @Input() public genders: Gender[] = [];
  @Input() public languages: string[] = ['de', 'en'];
  @Output() public submitData: EventEmitter<Profile.AsObject> = new EventEmitter<Profile.AsObject>();
  @Output() public changedLanguage: EventEmitter<string> = new EventEmitter<string>();
  @Output() public changeUsernameClicked: EventEmitter<void> = new EventEmitter();
  @Output() public avatarChanged: EventEmitter<void> = new EventEmitter();

  public profileForm!: UntypedFormGroup;

  private sub: Subscription = new Subscription();

  constructor(private fb: UntypedFormBuilder, private dialog: MatDialog) {
    this.profileForm = this.fb.group({
      userName: [{ value: '', disabled: true }, [requiredValidator]],
      firstName: [{ value: '', disabled: this.disabled }, requiredValidator],
      lastName: [{ value: '', disabled: this.disabled }, requiredValidator],
      nickName: [{ value: '', disabled: this.disabled }],
      displayName: [{ value: '', disabled: this.disabled }, requiredValidator],
      gender: [{ value: 0, disabled: this.disabled }],
      preferredLanguage: [{ value: '', disabled: this.disabled }],
    });
  }

  public ngOnChanges(): void {
    this.profileForm = this.fb.group({
      userName: [{ value: '', disabled: true }, [requiredValidator]],
      firstName: [{ value: '', disabled: this.disabled }, requiredValidator],
      lastName: [{ value: '', disabled: this.disabled }, requiredValidator],
      nickName: [{ value: '', disabled: this.disabled }],
      displayName: [{ value: '', disabled: this.disabled }, requiredValidator],
      gender: [{ value: 0, disabled: this.disabled }],
      preferredLanguage: [{ value: '', disabled: this.disabled }],
    });

    this.profileForm.patchValue({ userName: this.username, ...this.user.profile });

    if (this.preferredLanguage) {
      this.sub = this.preferredLanguage.valueChanges.subscribe((value) => {
        this.changedLanguage.emit(value);
      });
    }
  }

  public ngOnDestroy(): void {
    this.sub.unsubscribe();
  }

  public submitForm(): void {
    this.submitData.emit(this.profileForm.value);
  }

  public changeUsername(): void {
    this.changeUsernameClicked.emit();
  }

  public openUploadDialog(): void {
    const dialogRef = this.dialog.open(ProfilePictureComponent, {
      data: {
        profilePic: this.user.profile?.avatarUrl,
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((shouldReload) => {
      if (shouldReload) {
        this.avatarChanged.emit();
      }
    });
  }

  public get userName(): AbstractControl | null {
    return this.profileForm.get('userName');
  }

  public get firstName(): AbstractControl | null {
    return this.profileForm.get('firstName');
  }
  public get lastName(): AbstractControl | null {
    return this.profileForm.get('lastName');
  }
  public get nickName(): AbstractControl | null {
    return this.profileForm.get('nickName');
  }
  public get displayName(): AbstractControl | null {
    return this.profileForm.get('displayName');
  }
  public get gender(): AbstractControl | null {
    return this.profileForm.get('gender');
  }
  public get preferredLanguage(): AbstractControl | null {
    return this.profileForm.get('preferredLanguage');
  }
}
