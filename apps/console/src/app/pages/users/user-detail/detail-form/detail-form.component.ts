import { Component, DestroyRef, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { FormBuilder, FormControl } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { combineLatestWith, distinctUntilChanged, ReplaySubject } from 'rxjs';
import { requiredValidator } from 'src/app/modules/form-field/validators/validators';
import { ProfilePictureComponent } from './profile-picture/profile-picture.component';
import { Gender, HumanProfile, HumanProfileSchema } from '@zitadel/proto/zitadel/user/v2/user_pb';
import { filter, startWith } from 'rxjs/operators';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { Profile } from '@zitadel/proto/zitadel/user_pb';
//@ts-ignore
import { create } from '@zitadel/client';

function toHumanProfile(profile: HumanProfile | Profile): HumanProfile {
  if (profile.$typeName === 'zitadel.user.v2.HumanProfile') {
    return profile;
  }

  return create(HumanProfileSchema, {
    givenName: profile.firstName,
    familyName: profile.lastName,
    nickName: profile.nickName,
    displayName: profile.displayName,
    preferredLanguage: profile.preferredLanguage,
    gender: profile.gender,
    avatarUrl: profile.avatarUrl,
  });
}

@Component({
  selector: 'cnsl-detail-form',
  templateUrl: './detail-form.component.html',
  styleUrls: ['./detail-form.component.scss'],
})
export class DetailFormComponent implements OnInit {
  @Input() public showEditImage: boolean = false;
  @Input() public preferredLoginName: string = '';
  @Input({ required: true }) public set username(username: string) {
    this.username$.next(username);
  }
  @Input({ required: true }) public set profile(profile: HumanProfile | Profile) {
    this.profile$.next(toHumanProfile(profile));
  }
  @Input() public set disabled(disabled: boolean) {
    this.disabled$.next(disabled);
  }
  @Input() public genders: Gender[] = [];
  @Input() public languages: string[] = ['de', 'en'];
  @Output() public changedLanguage: EventEmitter<string> = new EventEmitter<string>();
  @Output() public changeUsernameClicked: EventEmitter<void> = new EventEmitter();
  @Output() public avatarChanged: EventEmitter<void> = new EventEmitter();

  private username$ = new ReplaySubject<string>(1);
  public profile$ = new ReplaySubject<HumanProfile>(1);
  public profileForm!: ReturnType<typeof this.buildForm>;
  public disabled$ = new ReplaySubject<boolean>(1);
  @Output() public submitData = new EventEmitter<HumanProfile>();

  constructor(
    private readonly fb: FormBuilder,
    private readonly dialog: MatDialog,
    private readonly destroyRef: DestroyRef,
  ) {
    this.profileForm = this.buildForm();
  }

  ngOnInit(): void {
    this.profileForm.controls.preferredLanguage.valueChanges
      .pipe(takeUntilDestroyed(this.destroyRef))
      .subscribe((value) => this.changedLanguage.emit(value));
  }

  private buildForm() {
    const form = this.fb.group({
      username: new FormControl('', { nonNullable: true, validators: [requiredValidator] }),
      givenName: new FormControl('', { nonNullable: true, validators: [requiredValidator] }),
      familyName: new FormControl('', { nonNullable: true, validators: [requiredValidator] }),
      nickName: new FormControl('', { nonNullable: true }),
      displayName: new FormControl('', { nonNullable: true, validators: [requiredValidator] }),
      preferredLanguage: new FormControl('', { nonNullable: true }),
      gender: new FormControl(Gender.UNSPECIFIED, { nonNullable: true }),
    });

    form.controls.username.disable();
    this.disabled$
      .pipe(startWith(true), distinctUntilChanged(), takeUntilDestroyed(this.destroyRef))
      .subscribe((disabled) => {
        this.toggleFormControl(form.controls.givenName, disabled);
        this.toggleFormControl(form.controls.familyName, disabled);
        this.toggleFormControl(form.controls.nickName, disabled);
        this.toggleFormControl(form.controls.displayName, disabled);
        this.toggleFormControl(form.controls.gender, disabled);
        this.toggleFormControl(form.controls.preferredLanguage, disabled);
      });

    this.username$
      .pipe(combineLatestWith(this.profile$), takeUntilDestroyed(this.destroyRef))
      .subscribe(([username, profile]) => {
        form.patchValue({
          username: username,
          ...profile,
        });
      });

    return form;
  }

  public submitForm(profile: HumanProfile): void {
    this.submitData.emit({ ...profile, ...this.profileForm.getRawValue() });
  }

  public changeUsername(): void {
    this.changeUsernameClicked.emit();
  }

  public openUploadDialog(profile: HumanProfile): void {
    const data = {
      profilePic: profile.avatarUrl,
    };

    const dialogRef = this.dialog.open<ProfilePictureComponent, typeof data, boolean>(ProfilePictureComponent, {
      width: '400px',
    });

    dialogRef
      .afterClosed()
      .pipe(filter(Boolean))
      .subscribe(() => {
        this.avatarChanged.emit();
      });
  }

  public toggleFormControl<T>(control: FormControl<T>, disabled: boolean) {
    if (disabled) {
      control.disable();
      return;
    }
    control.enable();
  }
}
