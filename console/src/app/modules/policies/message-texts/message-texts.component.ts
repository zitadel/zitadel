import { Component, Injector, Input, OnDestroy, OnInit, Type } from '@angular/core';
import { MatLegacyDialog as MatDialog } from '@angular/material/legacy-dialog';
import { MatLegacySelectChange as MatSelectChange } from '@angular/material/legacy-select';
import { BehaviorSubject, from, Observable, of, Subscription } from 'rxjs';
import {
  GetDefaultDomainClaimedMessageTextRequest as AdminGetDefaultDomainClaimedMessageTextRequest,
  GetDefaultInitMessageTextRequest as AdminGetDefaultInitMessageTextRequest,
  GetDefaultPasswordChangeMessageTextRequest as AdminGetDefaultPasswordChangeMessageTextRequest,
  GetDefaultPasswordlessRegistrationMessageTextRequest as AdminGetDefaultPasswordlessRegistrationMessageTextRequest,
  GetDefaultPasswordResetMessageTextRequest as AdminGetDefaultPasswordResetMessageTextRequest,
  GetDefaultVerifyEmailMessageTextRequest as AdminGetDefaultVerifyEmailMessageTextRequest,
  GetDefaultVerifyEmailOTPMessageTextRequest as AdminGetDefaultVerifyEmailOTPMessageTextRequest,
  GetDefaultVerifyPhoneMessageTextRequest as AdminGetDefaultVerifyPhoneMessageTextRequest,
  GetDefaultVerifySMSOTPMessageTextRequest as AdminGetDefaultVerifySMSOTPMessageTextRequest,
  SetDefaultDomainClaimedMessageTextRequest,
  SetDefaultInitMessageTextRequest,
  SetDefaultPasswordChangeMessageTextRequest,
  SetDefaultPasswordlessRegistrationMessageTextRequest,
  SetDefaultPasswordResetMessageTextRequest,
  SetDefaultVerifyEmailMessageTextRequest,
  SetDefaultVerifyEmailOTPMessageTextRequest,
  SetDefaultVerifyPhoneMessageTextRequest,
  SetDefaultVerifySMSOTPMessageTextRequest,
} from 'src/app/proto/generated/zitadel/admin_pb';
import {
  GetCustomDomainClaimedMessageTextRequest,
  GetCustomInitMessageTextRequest,
  GetCustomPasswordChangeMessageTextRequest,
  GetCustomPasswordlessRegistrationMessageTextRequest,
  GetCustomPasswordResetMessageTextRequest,
  GetCustomVerifyEmailMessageTextRequest,
  GetCustomVerifyEmailOTPMessageTextRequest,
  GetCustomVerifyPhoneMessageTextRequest,
  GetCustomVerifySMSOTPMessageTextRequest,
  GetDefaultDomainClaimedMessageTextRequest,
  GetDefaultInitMessageTextRequest,
  GetDefaultPasswordChangeMessageTextRequest,
  GetDefaultPasswordlessRegistrationMessageTextRequest,
  GetDefaultPasswordResetMessageTextRequest,
  GetDefaultVerifyEmailMessageTextRequest,
  GetDefaultVerifyEmailOTPMessageTextRequest,
  GetDefaultVerifyPhoneMessageTextRequest,
  GetDefaultVerifySMSOTPMessageTextRequest,
  SetCustomDomainClaimedMessageTextRequest,
  SetCustomInitMessageTextRequest,
  SetCustomPasswordChangeMessageTextRequest,
  SetCustomPasswordlessRegistrationMessageTextRequest,
  SetCustomPasswordResetMessageTextRequest,
  SetCustomVerifyEmailMessageTextRequest,
  SetCustomVerifyEmailOTPMessageTextRequest,
  SetCustomVerifyPhoneMessageTextRequest,
  SetCustomVerifySMSOTPMessageTextRequest,
} from 'src/app/proto/generated/zitadel/management_pb';
import { MessageCustomText } from 'src/app/proto/generated/zitadel/text_pb';
import { AdminService } from 'src/app/services/admin.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { supportedLanguages } from 'src/app/utils/language';
import { InfoSectionType } from '../../info-section/info-section.component';
import { WarnDialogComponent } from '../../warn-dialog/warn-dialog.component';
import { PolicyComponentServiceType } from '../policy-component-types.enum';

enum MESSAGETYPES {
  INIT = 'INIT',
  VERIFYPHONE = 'VP',
  VERIFYEMAIL = 'VE',
  PASSWORDRESET = 'PR',
  DOMAINCLAIMED = 'DC',
  PASSWORDLESS = 'PL',
  PASSWORDCHANGE = 'PC',
  VERIFYSMSOTP = 'VSO',
  VERIFYEMAILOTP = 'VEO',
}

const REQUESTMAP = {
  [PolicyComponentServiceType.MGMT]: {
    [MESSAGETYPES.PASSWORDCHANGE]: {
      get: new GetCustomPasswordChangeMessageTextRequest(),
      set: new SetCustomPasswordChangeMessageTextRequest(),
      getDefault: new GetDefaultPasswordChangeMessageTextRequest(),
      setFcn: (map: Partial<MessageCustomText.AsObject>): SetCustomPasswordChangeMessageTextRequest => {
        const req = new SetCustomPasswordChangeMessageTextRequest();
        req.setButtonText(map.buttonText ?? '');
        req.setFooterText(map.footerText ?? '');
        req.setGreeting(map.greeting ?? '');
        req.setPreHeader(map.preHeader ?? '');
        req.setSubject(map.subject ?? '');
        req.setText(map.text ?? '');
        req.setTitle(map.title ?? '');

        return req;
      },
    },
    [MESSAGETYPES.INIT]: {
      get: new GetCustomInitMessageTextRequest(),
      set: new SetCustomInitMessageTextRequest(),
      getDefault: new GetDefaultInitMessageTextRequest(),
      setFcn: (map: Partial<MessageCustomText.AsObject>): SetCustomInitMessageTextRequest => {
        const req = new SetCustomInitMessageTextRequest();
        req.setButtonText(map.buttonText ?? '');
        req.setFooterText(map.footerText ?? '');
        req.setGreeting(map.greeting ?? '');
        req.setPreHeader(map.preHeader ?? '');
        req.setSubject(map.subject ?? '');
        req.setText(map.text ?? '');
        req.setTitle(map.title ?? '');

        return req;
      },
    },
    [MESSAGETYPES.VERIFYEMAIL]: {
      get: new GetCustomVerifyEmailMessageTextRequest(),
      set: new SetCustomVerifyEmailMessageTextRequest(),
      getDefault: new GetDefaultVerifyEmailMessageTextRequest(),
      setFcn: (map: Partial<MessageCustomText.AsObject>): SetCustomVerifyEmailMessageTextRequest => {
        const req = new SetCustomVerifyEmailMessageTextRequest();
        req.setButtonText(map.buttonText ?? '');
        req.setFooterText(map.footerText ?? '');
        req.setGreeting(map.greeting ?? '');
        req.setPreHeader(map.preHeader ?? '');
        req.setSubject(map.subject ?? '');
        req.setText(map.text ?? '');
        req.setTitle(map.title ?? '');

        return req;
      },
    },
    [MESSAGETYPES.VERIFYPHONE]: {
      get: new GetCustomVerifyPhoneMessageTextRequest(),
      set: new SetCustomVerifyPhoneMessageTextRequest(),
      getDefault: new GetDefaultVerifyPhoneMessageTextRequest(),
      setFcn: (map: Partial<MessageCustomText.AsObject>): SetCustomVerifyPhoneMessageTextRequest => {
        const req = new SetCustomVerifyPhoneMessageTextRequest();
        req.setButtonText(map.buttonText ?? '');
        req.setFooterText(map.footerText ?? '');
        req.setGreeting(map.greeting ?? '');
        req.setPreHeader(map.preHeader ?? '');
        req.setSubject(map.subject ?? '');
        req.setText(map.text ?? '');
        req.setTitle(map.title ?? '');

        return req;
      },
    },
    [MESSAGETYPES.VERIFYSMSOTP]: {
      get: new GetCustomVerifySMSOTPMessageTextRequest(),
      set: new SetCustomVerifySMSOTPMessageTextRequest(),
      getDefault: new GetDefaultVerifySMSOTPMessageTextRequest(),
      setFcn: (map: Partial<MessageCustomText.AsObject>): SetCustomVerifySMSOTPMessageTextRequest => {
        const req = new SetCustomVerifySMSOTPMessageTextRequest();
        req.setText(map.text ?? '');

        return req;
      },
    },
    [MESSAGETYPES.VERIFYEMAILOTP]: {
      get: new GetCustomVerifyEmailOTPMessageTextRequest(),
      set: new SetCustomVerifyEmailOTPMessageTextRequest(),
      getDefault: new GetDefaultVerifyEmailOTPMessageTextRequest(),
      setFcn: (map: Partial<MessageCustomText.AsObject>): SetCustomVerifyEmailOTPMessageTextRequest => {
        const req = new SetCustomVerifyEmailOTPMessageTextRequest();
        req.setFooterText(map.footerText ?? '');
        req.setGreeting(map.greeting ?? '');
        req.setPreHeader(map.preHeader ?? '');
        req.setSubject(map.subject ?? '');
        req.setText(map.text ?? '');
        req.setTitle(map.title ?? '');

        return req;
      },
    },
    [MESSAGETYPES.PASSWORDRESET]: {
      get: new GetCustomPasswordResetMessageTextRequest(),
      set: new SetCustomPasswordResetMessageTextRequest(),
      getDefault: new GetDefaultPasswordResetMessageTextRequest(),
      setFcn: (
        map: Partial<SetCustomPasswordResetMessageTextRequest.AsObject>,
      ): SetCustomPasswordResetMessageTextRequest => {
        const req = new SetCustomPasswordResetMessageTextRequest();
        req.setButtonText(map.buttonText ?? '');
        req.setFooterText(map.footerText ?? '');
        req.setGreeting(map.greeting ?? '');
        req.setPreHeader(map.preHeader ?? '');
        req.setSubject(map.subject ?? '');
        req.setText(map.text ?? '');
        req.setTitle(map.title ?? '');

        return req;
      },
    },
    [MESSAGETYPES.DOMAINCLAIMED]: {
      get: new GetCustomDomainClaimedMessageTextRequest(),
      set: new SetCustomDomainClaimedMessageTextRequest(),
      getDefault: new GetDefaultDomainClaimedMessageTextRequest(),
      setFcn: (
        map: Partial<SetCustomDomainClaimedMessageTextRequest.AsObject>,
      ): SetCustomDomainClaimedMessageTextRequest => {
        const req = new SetCustomDomainClaimedMessageTextRequest();
        req.setButtonText(map.buttonText ?? '');
        req.setFooterText(map.footerText ?? '');
        req.setGreeting(map.greeting ?? '');
        req.setPreHeader(map.preHeader ?? '');
        req.setSubject(map.subject ?? '');
        req.setText(map.text ?? '');
        req.setTitle(map.title ?? '');

        return req;
      },
    },
    [MESSAGETYPES.PASSWORDLESS]: {
      get: new GetCustomPasswordlessRegistrationMessageTextRequest(),
      set: new SetCustomPasswordlessRegistrationMessageTextRequest(),
      getDefault: new GetDefaultPasswordlessRegistrationMessageTextRequest(),
      setFcn: (
        map: Partial<SetCustomPasswordlessRegistrationMessageTextRequest.AsObject>,
      ): SetCustomPasswordlessRegistrationMessageTextRequest => {
        const req = new SetCustomPasswordlessRegistrationMessageTextRequest();
        req.setButtonText(map.buttonText ?? '');
        req.setFooterText(map.footerText ?? '');
        req.setGreeting(map.greeting ?? '');
        req.setPreHeader(map.preHeader ?? '');
        req.setSubject(map.subject ?? '');
        req.setText(map.text ?? '');
        req.setTitle(map.title ?? '');

        return req;
      },
    },
  },
  [PolicyComponentServiceType.ADMIN]: {
    [MESSAGETYPES.PASSWORDCHANGE]: {
      get: new AdminGetDefaultPasswordChangeMessageTextRequest(),
      set: new SetDefaultPasswordChangeMessageTextRequest(),
      setFcn: (map: Partial<MessageCustomText.AsObject>): SetDefaultPasswordChangeMessageTextRequest => {
        const req = new SetDefaultPasswordChangeMessageTextRequest();
        req.setButtonText(map.buttonText ?? '');
        req.setFooterText(map.footerText ?? '');
        req.setGreeting(map.greeting ?? '');
        req.setPreHeader(map.preHeader ?? '');
        req.setSubject(map.subject ?? '');
        req.setText(map.text ?? '');
        req.setTitle(map.title ?? '');

        return req;
      },
    },
    [MESSAGETYPES.INIT]: {
      get: new AdminGetDefaultInitMessageTextRequest(),
      set: new SetDefaultInitMessageTextRequest(),
      setFcn: (map: Partial<MessageCustomText.AsObject>): SetDefaultInitMessageTextRequest => {
        const req = new SetDefaultInitMessageTextRequest();
        req.setButtonText(map.buttonText ?? '');
        req.setFooterText(map.footerText ?? '');
        req.setGreeting(map.greeting ?? '');
        req.setPreHeader(map.preHeader ?? '');
        req.setSubject(map.subject ?? '');
        req.setText(map.text ?? '');
        req.setTitle(map.title ?? '');

        return req;
      },
    },
    [MESSAGETYPES.VERIFYEMAIL]: {
      get: new AdminGetDefaultVerifyEmailMessageTextRequest(),
      set: new SetDefaultVerifyEmailMessageTextRequest(),
      setFcn: (map: Partial<MessageCustomText.AsObject>): SetDefaultVerifyEmailMessageTextRequest => {
        const req = new SetDefaultVerifyEmailMessageTextRequest();
        req.setButtonText(map.buttonText ?? '');
        req.setFooterText(map.footerText ?? '');
        req.setGreeting(map.greeting ?? '');
        req.setPreHeader(map.preHeader ?? '');
        req.setSubject(map.subject ?? '');
        req.setText(map.text ?? '');
        req.setTitle(map.title ?? '');

        return req;
      },
    },
    [MESSAGETYPES.VERIFYPHONE]: {
      get: new AdminGetDefaultVerifyPhoneMessageTextRequest(),
      set: new SetDefaultVerifyPhoneMessageTextRequest(),
      setFcn: (map: Partial<MessageCustomText.AsObject>): SetDefaultVerifyPhoneMessageTextRequest => {
        const req = new SetDefaultVerifyPhoneMessageTextRequest();
        req.setButtonText(map.buttonText ?? '');
        req.setFooterText(map.footerText ?? '');
        req.setGreeting(map.greeting ?? '');
        req.setPreHeader(map.preHeader ?? '');
        req.setSubject(map.subject ?? '');
        req.setText(map.text ?? '');
        req.setTitle(map.title ?? '');

        return req;
      },
    },
    [MESSAGETYPES.VERIFYSMSOTP]: {
      get: new AdminGetDefaultVerifySMSOTPMessageTextRequest(),
      set: new SetDefaultVerifySMSOTPMessageTextRequest(),
      setFcn: (map: Partial<MessageCustomText.AsObject>): SetDefaultVerifySMSOTPMessageTextRequest => {
        const req = new SetDefaultVerifySMSOTPMessageTextRequest();
        req.setText(map.text ?? '');

        return req;
      },
    },
    [MESSAGETYPES.VERIFYEMAILOTP]: {
      get: new AdminGetDefaultVerifyEmailOTPMessageTextRequest(),
      set: new SetDefaultVerifyEmailOTPMessageTextRequest(),
      setFcn: (map: Partial<MessageCustomText.AsObject>): SetDefaultVerifyEmailOTPMessageTextRequest => {
        const req = new SetDefaultVerifyEmailOTPMessageTextRequest();
        req.setFooterText(map.footerText ?? '');
        req.setGreeting(map.greeting ?? '');
        req.setPreHeader(map.preHeader ?? '');
        req.setSubject(map.subject ?? '');
        req.setText(map.text ?? '');
        req.setTitle(map.title ?? '');

        return req;
      },
    },
    [MESSAGETYPES.PASSWORDRESET]: {
      get: new AdminGetDefaultPasswordResetMessageTextRequest(),
      set: new SetDefaultPasswordResetMessageTextRequest(),
      setFcn: (
        map: Partial<SetDefaultPasswordResetMessageTextRequest.AsObject>,
      ): SetDefaultPasswordResetMessageTextRequest => {
        const req = new SetDefaultPasswordResetMessageTextRequest();
        req.setButtonText(map.buttonText ?? '');
        req.setFooterText(map.footerText ?? '');
        req.setGreeting(map.greeting ?? '');
        req.setPreHeader(map.preHeader ?? '');
        req.setSubject(map.subject ?? '');
        req.setText(map.text ?? '');
        req.setTitle(map.title ?? '');

        return req;
      },
    },
    [MESSAGETYPES.DOMAINCLAIMED]: {
      get: new AdminGetDefaultDomainClaimedMessageTextRequest(),
      set: new SetDefaultDomainClaimedMessageTextRequest(),
      setFcn: (
        map: Partial<SetDefaultDomainClaimedMessageTextRequest.AsObject>,
      ): SetDefaultDomainClaimedMessageTextRequest => {
        const req = new SetDefaultDomainClaimedMessageTextRequest();
        req.setButtonText(map.buttonText ?? '');
        req.setFooterText(map.footerText ?? '');
        req.setGreeting(map.greeting ?? '');
        req.setPreHeader(map.preHeader ?? '');
        req.setSubject(map.subject ?? '');
        req.setText(map.text ?? '');
        req.setTitle(map.title ?? '');

        return req;
      },
    },
    [MESSAGETYPES.PASSWORDLESS]: {
      get: new AdminGetDefaultPasswordlessRegistrationMessageTextRequest(),
      set: new SetDefaultPasswordlessRegistrationMessageTextRequest(),
      setFcn: (
        map: Partial<SetDefaultPasswordlessRegistrationMessageTextRequest.AsObject>,
      ): SetDefaultPasswordlessRegistrationMessageTextRequest => {
        const req = new SetDefaultPasswordlessRegistrationMessageTextRequest();
        req.setButtonText(map.buttonText ?? '');
        req.setFooterText(map.footerText ?? '');
        req.setGreeting(map.greeting ?? '');
        req.setPreHeader(map.preHeader ?? '');
        req.setSubject(map.subject ?? '');
        req.setText(map.text ?? '');
        req.setTitle(map.title ?? '');

        return req;
      },
    },
  },
};
@Component({
  selector: 'cnsl-message-texts',
  templateUrl: './message-texts.component.html',
  styleUrls: ['./message-texts.component.scss'],
})
export class MessageTextsComponent implements OnInit, OnDestroy {
  public loading: boolean = false;
  public getDefaultMessageTextMap$: Observable<{ [key: string]: string }> = of({});
  public getCustomMessageTextMap$: BehaviorSubject<{ [key: string]: string | boolean }> = new BehaviorSubject({}); // boolean because of isDefault

  public currentType: MESSAGETYPES = MESSAGETYPES.INIT;

  public service!: ManagementService | AdminService;
  public PolicyComponentServiceType: any = PolicyComponentServiceType;
  @Input() public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;

  public MESSAGETYPES: any = MESSAGETYPES;

  public updateRequest!: any;

  public InfoSectionType: any = InfoSectionType;
  public chips: {
    [messagetype: string]: Array<{ key: string; value: string }>;
  } = {
    [MESSAGETYPES.DOMAINCLAIMED]: [
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.preferredLoginName', value: '{{.PreferredLoginName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.domain', value: '{{.Domain}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.tempUsername', value: '{{.TempUsername}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.username', value: '{{.UserName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.firstname', value: '{{.FirstName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.lastname', value: '{{.Lastname}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.nickName', value: '{{.NickName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.displayName', value: '{{.DisplayName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.lastEmail', value: '{{.LastEmail}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.verifiedEmail', value: '{{.VerifiedEmail}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.lastPhone', value: '{{.LastPhone}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.verifiedPhone', value: '{{.VerifiedPhone}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.loginnames', value: '{{.LoginNames}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.changedate', value: '{{.ChangeDate}}' },
    ],
    [MESSAGETYPES.INIT]: [
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.code', value: '{{.Code}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.preferredLoginName', value: '{{.PreferredLoginName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.username', value: '{{.UserName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.firstname', value: '{{.FirstName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.lastname', value: '{{.Lastname}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.nickName', value: '{{.NickName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.displayName', value: '{{.DisplayName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.lastEmail', value: '{{.LastEmail}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.verifiedEmail', value: '{{.VerifiedEmail}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.lastPhone', value: '{{.LastPhone}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.verifiedPhone', value: '{{.VerifiedPhone}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.loginnames', value: '{{.LoginNames}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.changedate', value: '{{.ChangeDate}}' },
    ],
    [MESSAGETYPES.PASSWORDRESET]: [
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.code', value: '{{.Code}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.preferredLoginName', value: '{{.PreferredLoginName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.username', value: '{{.UserName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.firstname', value: '{{.FirstName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.lastname', value: '{{.Lastname}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.nickName', value: '{{.NickName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.displayName', value: '{{.DisplayName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.lastEmail', value: '{{.LastEmail}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.verifiedEmail', value: '{{.VerifiedEmail}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.lastPhone', value: '{{.LastPhone}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.verifiedPhone', value: '{{.VerifiedPhone}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.loginnames', value: '{{.LoginNames}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.changedate', value: '{{.ChangeDate}}' },
    ],
    [MESSAGETYPES.VERIFYEMAIL]: [
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.code', value: '{{.Code}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.preferredLoginName', value: '{{.PreferredLoginName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.username', value: '{{.UserName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.firstname', value: '{{.FirstName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.lastname', value: '{{.Lastname}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.nickName', value: '{{.NickName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.displayName', value: '{{.DisplayName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.lastEmail', value: '{{.LastEmail}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.verifiedEmail', value: '{{.VerifiedEmail}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.lastPhone', value: '{{.LastPhone}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.verifiedPhone', value: '{{.VerifiedPhone}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.loginnames', value: '{{.LoginNames}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.changedate', value: '{{.ChangeDate}}' },
    ],
    [MESSAGETYPES.VERIFYPHONE]: [
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.code', value: '{{.Code}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.preferredLoginName', value: '{{.PreferredLoginName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.username', value: '{{.UserName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.firstname', value: '{{.FirstName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.lastname', value: '{{.Lastname}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.nickName', value: '{{.NickName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.displayName', value: '{{.DisplayName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.lastEmail', value: '{{.LastEmail}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.verifiedEmail', value: '{{.VerifiedEmail}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.lastPhone', value: '{{.LastPhone}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.verifiedPhone', value: '{{.VerifiedPhone}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.loginnames', value: '{{.LoginNames}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.changedate', value: '{{.ChangeDate}}' },
    ],
    [MESSAGETYPES.VERIFYSMSOTP]: [
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.otp', value: '{{.OTP}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.verifyUrl', value: '{{.VerifyURL}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.preferredLoginName', value: '{{.PreferredLoginName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.username', value: '{{.UserName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.firstname', value: '{{.FirstName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.lastname', value: '{{.Lastname}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.nickName', value: '{{.NickName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.displayName', value: '{{.DisplayName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.lastEmail', value: '{{.LastEmail}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.verifiedEmail', value: '{{.VerifiedEmail}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.lastPhone', value: '{{.LastPhone}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.verifiedPhone', value: '{{.VerifiedPhone}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.loginnames', value: '{{.LoginNames}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.changedate', value: '{{.ChangeDate}}' },
    ],
    [MESSAGETYPES.VERIFYEMAILOTP]: [
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.otp', value: '{{.OTP}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.verifyUrl', value: '{{.VerifyURL}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.preferredLoginName', value: '{{.PreferredLoginName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.username', value: '{{.UserName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.firstname', value: '{{.FirstName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.lastname', value: '{{.Lastname}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.nickName', value: '{{.NickName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.displayName', value: '{{.DisplayName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.lastEmail', value: '{{.LastEmail}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.verifiedEmail', value: '{{.VerifiedEmail}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.lastPhone', value: '{{.LastPhone}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.verifiedPhone', value: '{{.VerifiedPhone}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.loginnames', value: '{{.LoginNames}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.changedate', value: '{{.ChangeDate}}' },
    ],
    [MESSAGETYPES.PASSWORDLESS]: [
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.preferredLoginName', value: '{{.PreferredLoginName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.username', value: '{{.UserName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.firstname', value: '{{.FirstName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.lastname', value: '{{.Lastname}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.nickName', value: '{{.NickName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.displayName', value: '{{.DisplayName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.lastEmail', value: '{{.LastEmail}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.verifiedEmail', value: '{{.VerifiedEmail}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.lastPhone', value: '{{.LastPhone}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.verifiedPhone', value: '{{.VerifiedPhone}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.loginnames', value: '{{.LoginNames}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.changedate', value: '{{.ChangeDate}}' },
    ],
    [MESSAGETYPES.PASSWORDCHANGE]: [
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.preferredLoginName', value: '{{.PreferredLoginName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.username', value: '{{.UserName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.firstname', value: '{{.FirstName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.lastname', value: '{{.Lastname}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.nickName', value: '{{.NickName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.displayName', value: '{{.DisplayName}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.lastEmail', value: '{{.LastEmail}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.verifiedEmail', value: '{{.VerifiedEmail}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.lastPhone', value: '{{.LastPhone}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.verifiedPhone', value: '{{.VerifiedPhone}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.loginnames', value: '{{.LoginNames}}' },
      { key: 'POLICY.MESSAGE_TEXTS.CHIPS.changedate', value: '{{.ChangeDate}}' },
    ],
  };

  public locale: string = 'en';
  public LOCALES: string[] = supportedLanguages;
  private sub: Subscription = new Subscription();
  public canWrite$: Observable<boolean> = this.authService.isAllowed([
    this.serviceType === PolicyComponentServiceType.ADMIN
      ? 'iam.policy.write'
      : this.serviceType === PolicyComponentServiceType.MGMT
      ? 'policy.write'
      : '',
  ]);

  constructor(
    private authService: GrpcAuthService,
    private toast: ToastService,
    private injector: Injector,
    private dialog: MatDialog,
  ) {}

  ngOnInit(): void {
    switch (this.serviceType) {
      case PolicyComponentServiceType.MGMT:
        this.service = this.injector.get(ManagementService as Type<ManagementService>);
        this.service.getSupportedLanguages().then((lang) => {
          this.LOCALES = lang.languagesList;
        });
        this.loadData(this.currentType);
        break;
      case PolicyComponentServiceType.ADMIN:
        this.service = this.injector.get(AdminService as Type<AdminService>);
        this.service.getSupportedLanguages().then((lang) => {
          this.LOCALES = lang.languagesList;
        });
        this.loadData(this.currentType);
        break;
    }
  }

  public getDefaultValues(type: MESSAGETYPES, req: any): Promise<any> {
    switch (type) {
      case MESSAGETYPES.INIT:
        return this.stripEmail(this.service.getDefaultInitMessageText(req));
      case MESSAGETYPES.VERIFYEMAIL:
        return this.stripEmail(this.service.getDefaultVerifyEmailMessageText(req));
      case MESSAGETYPES.VERIFYPHONE:
        return this.stripSMS(this.service.getDefaultVerifyPhoneMessageText(req));
      case MESSAGETYPES.VERIFYSMSOTP:
        return this.stripSMS(this.service.getDefaultVerifySMSOTPMessageText(req));
      case MESSAGETYPES.VERIFYEMAILOTP:
        return this.stripEmail(this.service.getDefaultVerifyEmailOTPMessageText(req));
      case MESSAGETYPES.PASSWORDRESET:
        return this.stripEmail(this.service.getDefaultPasswordResetMessageText(req));
      case MESSAGETYPES.DOMAINCLAIMED:
        return this.stripEmail(this.service.getDefaultDomainClaimedMessageText(req));
      case MESSAGETYPES.PASSWORDLESS:
        return this.stripEmail(this.service.getDefaultPasswordlessRegistrationMessageText(req));
      case MESSAGETYPES.PASSWORDCHANGE:
        return this.stripEmail(this.service.getDefaultPasswordChangeMessageText(req));
    }
  }

  public getCurrentValues(type: MESSAGETYPES, req: any): Promise<any> | undefined {
    switch (type) {
      case MESSAGETYPES.INIT:
        return this.stripEmail(this.service.getCustomInitMessageText(req));
      case MESSAGETYPES.VERIFYEMAIL:
        return this.stripEmail(this.service.getCustomVerifyEmailMessageText(req));
      case MESSAGETYPES.VERIFYPHONE:
        return this.stripSMS(this.service.getCustomVerifyPhoneMessageText(req));
      case MESSAGETYPES.VERIFYSMSOTP:
        return this.stripSMS(this.service.getCustomVerifySMSOTPMessageText(req));
      case MESSAGETYPES.VERIFYEMAILOTP:
        return this.stripEmail(this.service.getCustomVerifyEmailOTPMessageText(req));
      case MESSAGETYPES.PASSWORDRESET:
        return this.stripEmail(this.service.getCustomPasswordResetMessageText(req));
      case MESSAGETYPES.DOMAINCLAIMED:
        return this.stripEmail(this.service.getCustomDomainClaimedMessageText(req));
      case MESSAGETYPES.PASSWORDLESS:
        return this.stripEmail(this.service.getCustomPasswordlessRegistrationMessageText(req));
      case MESSAGETYPES.PASSWORDCHANGE:
        return this.stripEmail(this.service.getCustomPasswordChangeMessageText(req));
      default:
        return undefined;
    }
  }

  public changeLocale(selection: MatSelectChange): void {
    this.locale = selection.value;
    this.loadData(this.currentType);
  }

  public async loadData(type: MESSAGETYPES): Promise<any> {
    if (this.serviceType === PolicyComponentServiceType.MGMT) {
      const reqDefaultInit = REQUESTMAP[this.serviceType][type].getDefault;

      reqDefaultInit.setLanguage(this.locale);
      this.getDefaultMessageTextMap$ = from(this.getDefaultValues(type, reqDefaultInit));
    }

    const reqCustomInit = REQUESTMAP[this.serviceType][type].get.setLanguage(this.locale);
    this.loading = true;
    return this.getCurrentValues(type, reqCustomInit)
      ?.then((data) => {
        this.loading = false;
        this.getCustomMessageTextMap$.next(data);
      })
      .catch((error) => {
        this.loading = false;
        this.toast.showError(error);
      });
  }

  public updateCurrentValues(values: { [key: string]: string }): void {
    const req = REQUESTMAP[this.serviceType][this.currentType].setFcn;
    const mappedValues = req(values);
    this.updateRequest = mappedValues;
    this.updateRequest.setLanguage(this.locale);
  }

  public saveCurrentMessage(): any {
    const handler = (prom: Promise<any>): Promise<any> => {
      return prom
        .then(() => {
          this.toast.showInfo('POLICY.MESSAGE_TEXTS.TOAST.UPDATED', true);
        })
        .catch((error) => this.toast.showError(error));
    };
    if (this.serviceType === PolicyComponentServiceType.MGMT) {
      switch (this.currentType) {
        case MESSAGETYPES.INIT:
          return handler((this.service as ManagementService).setCustomInitMessageText(this.updateRequest));
        case MESSAGETYPES.VERIFYEMAIL:
          return handler((this.service as ManagementService).setCustomVerifyEmailMessageText(this.updateRequest));
        case MESSAGETYPES.VERIFYPHONE:
          return handler((this.service as ManagementService).setCustomVerifyPhoneMessageText(this.updateRequest));
        case MESSAGETYPES.VERIFYSMSOTP:
          return handler((this.service as ManagementService).setCustomVerifySMSOTPMessageText(this.updateRequest));
        case MESSAGETYPES.VERIFYEMAILOTP:
          return handler((this.service as ManagementService).setCustomVerifyEmailOTPMessageText(this.updateRequest));
        case MESSAGETYPES.PASSWORDRESET:
          return handler((this.service as ManagementService).setCustomPasswordResetMessageText(this.updateRequest));
        case MESSAGETYPES.DOMAINCLAIMED:
          return handler((this.service as ManagementService).setCustomDomainClaimedMessageCustomText(this.updateRequest));
        case MESSAGETYPES.PASSWORDLESS:
          return handler(
            (this.service as ManagementService).setCustomPasswordlessRegistrationMessageCustomText(this.updateRequest),
          );
        case MESSAGETYPES.PASSWORDCHANGE:
          return handler((this.service as ManagementService).setCustomPasswordChangeMessageText(this.updateRequest));
      }
    } else if (this.serviceType === PolicyComponentServiceType.ADMIN) {
      switch (this.currentType) {
        case MESSAGETYPES.INIT:
          return handler((this.service as AdminService).setDefaultInitMessageText(this.updateRequest));
        case MESSAGETYPES.VERIFYPHONE:
          return handler((this.service as AdminService).setDefaultVerifyPhoneMessageText(this.updateRequest));
        case MESSAGETYPES.VERIFYSMSOTP:
          return handler((this.service as AdminService).setDefaultVerifySMSOTPMessageText(this.updateRequest));
        case MESSAGETYPES.VERIFYEMAILOTP:
          return handler((this.service as AdminService).setDefaultVerifyEmailOTPMessageText(this.updateRequest));
        case MESSAGETYPES.VERIFYEMAIL:
          return handler((this.service as AdminService).setDefaultVerifyEmailMessageText(this.updateRequest));
        case MESSAGETYPES.PASSWORDRESET:
          return handler((this.service as AdminService).setDefaultPasswordResetMessageText(this.updateRequest));
        case MESSAGETYPES.DOMAINCLAIMED:
          return handler((this.service as AdminService).setDefaultDomainClaimedMessageText(this.updateRequest));
        case MESSAGETYPES.PASSWORDLESS:
          return handler((this.service as AdminService).setDefaultPasswordlessRegistrationMessageText(this.updateRequest));
        case MESSAGETYPES.PASSWORDCHANGE:
          return handler((this.service as AdminService).setDefaultPasswordChangeMessageText(this.updateRequest));
      }
    }
  }

  public resetDefault(): void {
    const dialogRef = this.dialog.open(WarnDialogComponent, {
      data: {
        icon: 'las la-history',
        confirmKey: 'ACTIONS.RESTORE',
        cancelKey: 'ACTIONS.CANCEL',
        titleKey: 'POLICY.LOGIN_TEXTS.RESET_TITLE',
        descriptionKey: 'POLICY.LOGIN_TEXTS.RESET_DESCRIPTION',
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (!resp) {
        return Promise.reject();
      }

      const handler = (prom: Promise<any>): Promise<any> => {
        return prom
          .then(() => {
            setTimeout(() => {
              this.loadData(this.currentType);
            }, 1000);
          })
          .catch((error) => {
            this.toast.showError(error);
          });
      };

      switch (this.currentType) {
        case MESSAGETYPES.INIT:
          return handler(this.service.resetCustomInitMessageTextToDefault(this.locale));
        case MESSAGETYPES.VERIFYPHONE:
          return handler(this.service.resetCustomVerifyPhoneMessageTextToDefault(this.locale));
        case MESSAGETYPES.VERIFYSMSOTP:
          return handler(this.service.resetCustomVerifySMSOTPMessageTextToDefault(this.locale));
        case MESSAGETYPES.VERIFYEMAILOTP:
          return handler(this.service.resetCustomVerifyEmailOTPMessageTextToDefault(this.locale));
        case MESSAGETYPES.VERIFYEMAIL:
          return handler(this.service.resetCustomVerifyEmailMessageTextToDefault(this.locale));
        case MESSAGETYPES.PASSWORDRESET:
          return handler(this.service.resetCustomPasswordResetMessageTextToDefault(this.locale));
        case MESSAGETYPES.DOMAINCLAIMED:
          return handler(this.service.resetCustomDomainClaimedMessageTextToDefault(this.locale));
        case MESSAGETYPES.PASSWORDLESS:
          return handler(this.service.resetCustomPasswordlessRegistrationMessageTextToDefault(this.locale));
        case MESSAGETYPES.PASSWORDCHANGE:
          return handler(this.service.resetCustomPasswordChangeMessageTextToDefault(this.locale));
        default:
          return Promise.reject();
      }
    });
  }

  private strip(prom: Promise<any>, properties: Array<string>): Promise<any> {
    return prom.then((res) => {
      if (res.customText) {
        properties.forEach((property) => {
          delete res.customText[property];
        });
        return Object.assign({}, res.customText as unknown as { [key: string]: string });
      } else {
        return {};
      }
    });
  }

  private stripEmail(prom: Promise<any>): Promise<any> {
    return this.strip(prom, ['details']);
  }

  private stripSMS(prom: Promise<any>): Promise<any> {
    return this.strip(prom, ['details', 'buttonText', 'footerText', 'greeting', 'preHeader', 'subject', 'title']);
  }

  public ngOnDestroy(): void {
    this.sub.unsubscribe();
  }

  public changedCurrentType(): void {
    this.loadData(this.currentType);
  }
}
