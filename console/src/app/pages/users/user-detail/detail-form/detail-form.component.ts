import { Component, EventEmitter, Input, OnChanges, OnDestroy, Output } from '@angular/core';
import { AbstractControl, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { Subscription } from 'rxjs';
import { Gender, Human, User } from 'src/app/proto/generated/zitadel/user_pb';

import { ProfilePictureComponent } from './profile-picture/profile-picture.component';

@Component({
  selector: 'app-detail-form',
  templateUrl: './detail-form.component.html',
  styleUrls: ['./detail-form.component.scss'],
})
export class DetailFormComponent implements OnDestroy, OnChanges {
  @Input() public showEditImage: boolean = false;
  @Input() public preferredLoginName: string = '';
  @Input() public username!: string;
  @Input() public user!: Human.AsObject;
  @Input() public disabled: boolean = false;
  @Input() public genders: Gender[] = [];
  @Input() public languages: string[] = ['de', 'en'];
  @Output() public submitData: EventEmitter<User> = new EventEmitter<User>();
  @Output() public changedLanguage: EventEmitter<string> = new EventEmitter<string>();

  public profileForm!: FormGroup;

  private sub: Subscription = new Subscription();

  constructor(
    private fb: FormBuilder,
    private dialog: MatDialog,
  ) {
    this.profileForm = this.fb.group({
      userName: [{ value: '', disabled: true }, [
        Validators.required,
      ]],
      firstName: [{ value: '', disabled: this.disabled }, Validators.required],
      lastName: [{ value: '', disabled: this.disabled }, Validators.required],
      nickName: [{ value: '', disabled: this.disabled }],
      displayName: [{ value: '', disabled: this.disabled }, Validators.required],
      gender: [{ value: 0, disabled: this.disabled }],
      preferredLanguage: [{ value: '', disabled: this.disabled }],
    });
  }

  public ngOnChanges(): void {
    this.profileForm = this.fb.group({
      userName: [{ value: '', disabled: true }, [
        Validators.required,
      ]],
      firstName: [{ value: '', disabled: this.disabled }, Validators.required],
      lastName: [{ value: '', disabled: this.disabled }, Validators.required],
      nickName: [{ value: '', disabled: this.disabled }],
      displayName: [{ value: '', disabled: this.disabled }, Validators.required],
      gender: [{ value: 0, disabled: this.disabled }],
      preferredLanguage: [{ value: '', disabled: this.disabled }],
    });

    this.profileForm.patchValue({ userName: this.username, ...this.user.profile });

    if (this.preferredLanguage) {
      this.sub = this.preferredLanguage.valueChanges.subscribe(value => {
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

  public openUploadDialog(): void {
    const dialogRef = this.dialog.open(ProfilePictureComponent, {
      data: {
        profilePic: this.user.profile?.avatarUrl,
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe(resp => {
      if (resp) {
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
