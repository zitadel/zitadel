import { OverlayModule } from '@angular/cdk/overlay';
import { CommonModule, registerLocaleData } from '@angular/common';
import { HttpClientModule } from '@angular/common/http';
import localeDe from '@angular/common/locales/de';
import { APP_INITIALIZER, NgModule } from '@angular/core';
import { ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatNativeDateModule } from '@angular/material/core';
import { MatDialogModule } from '@angular/material/dialog';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatTooltipModule } from '@angular/material/tooltip';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { ServiceWorkerModule } from '@angular/service-worker';
import { TranslateLoader, TranslateModule } from '@ngx-translate/core';
import { AuthConfig, OAuthModule, OAuthStorage } from 'angular-oauth2-oidc';
import { QuicklinkModule } from 'ngx-quicklink';
import { from, Observable } from 'rxjs';
import { OnboardingModule } from 'src/app/modules/onboarding/onboarding.module';
import { RegExpPipeModule } from 'src/app/pipes/regexp-pipe/regexp-pipe.module';
import { SubscriptionService } from 'src/app/services/subscription.service';
import { UploadService } from 'src/app/services/upload.service';

import { environment } from '../environments/environment';
import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { HasRoleModule } from './directives/has-role/has-role.module';
import { OutsideClickModule } from './directives/outside-click/outside-click.module';
import { AccountsCardModule } from './modules/accounts-card/accounts-card.module';
import { AvatarModule } from './modules/avatar/avatar.module';
import { InputModule } from './modules/input/input.module';
import { WarnDialogModule } from './modules/warn-dialog/warn-dialog.module';
import { SignedoutComponent } from './pages/signedout/signedout.component';
import { HasFeaturePipeModule } from './pipes/has-feature-pipe/has-feature-pipe.module';
import { HasRolePipeModule } from './pipes/has-role-pipe/has-role-pipe.module';
import { AuthenticationService } from './services/authentication.service';
import { GrpcAuthService } from './services/grpc-auth.service';
import { GrpcService } from './services/grpc.service';
import { AuthInterceptor } from './services/interceptors/auth.interceptor';
import { GRPC_INTERCEPTORS } from './services/interceptors/grpc-interceptor';
import { I18nInterceptor } from './services/interceptors/i18n.interceptor';
import { OrgInterceptor } from './services/interceptors/org.interceptor';
import { RefreshService } from './services/refresh.service';
import { SeoService } from './services/seo.service';
import { StatehandlerProcessorService, StatehandlerProcessorServiceImpl } from './services/statehandler-processor.service';
import { StatehandlerService, StatehandlerServiceImpl } from './services/statehandler.service';
import { StorageService } from './services/storage.service';
import { ThemeService } from './services/theme.service';

registerLocaleData(localeDe);

export class WebpackTranslateLoader implements TranslateLoader {
  getTranslation(lang: string): Observable<any> {
    return from(import(`../assets/i18n/${lang}.json`));
  }
}

const appInitializerFn = (grpcServ: GrpcService) => {
  return () => {
    return grpcServ.loadAppEnvironment();
  };
};

const stateHandlerFn = (stateHandler: StatehandlerService) => {
  return () => {
    return stateHandler.initStateHandler();
  };
};

const authConfig: AuthConfig = {
  scope: 'openid profile email', // offline_access
  responseType: 'code',
  oidc: true,
};

@NgModule({
  declarations: [
    AppComponent,
    SignedoutComponent,
  ],
  imports: [
    AppRoutingModule,
    CommonModule,
    BrowserModule,
    OverlayModule,
    OAuthModule.forRoot({
      resourceServer: {
        allowedUrls: ['https://test.api.zitadel.caos.ch/caos.zitadel.auth.api.v1.AuthService', 'https://test.api.zitadel.caos.ch/oauth/v2/userinfo', 'https://test.api.zitadel.caos.ch/caos.zitadel.management.api.v1.ManagementService/', 'https://preview.api.zitadel.caos.ch'],
        sendAccessToken: true,
      },
    }),
    TranslateModule.forRoot({
      loader: {
        provide: TranslateLoader,
        useClass: WebpackTranslateLoader,
      },
    }),
    MatNativeDateModule,
    QuicklinkModule,
    AccountsCardModule,
    HasRoleModule,
    BrowserAnimationsModule,
    HttpClientModule,
    MatButtonModule,
    MatIconModule,
    MatTooltipModule,
    MatSidenavModule,
    MatCardModule,
    OutsideClickModule,
    InputModule,
    HasRolePipeModule,
    HasFeaturePipeModule,
    MatProgressBarModule,
    MatProgressSpinnerModule,
    MatToolbarModule,
    ReactiveFormsModule,
    MatMenuModule,
    MatSnackBarModule,
    AvatarModule,
    WarnDialogModule,
    MatDialogModule,
    RegExpPipeModule,
    OnboardingModule,
    ServiceWorkerModule.register('ngsw-worker.js', { enabled: environment.production }),
  ],
  providers: [
    ThemeService,
    {
      provide: APP_INITIALIZER,
      useFactory: appInitializerFn,
      multi: true,
      deps: [GrpcService],
    },
    {
      provide: APP_INITIALIZER,
      useFactory: stateHandlerFn,
      multi: true,
      deps: [StatehandlerService],
    },
    {
      provide: AuthConfig,
      useValue: authConfig,
    },
    {
      provide: StatehandlerProcessorService,
      useClass: StatehandlerProcessorServiceImpl,
    },
    {
      provide: StatehandlerService,
      useClass: StatehandlerServiceImpl,
    },
    {
      provide: OAuthStorage,
      useClass: StorageService,
    },
    {
      provide: GRPC_INTERCEPTORS,
      multi: true,
      useClass: AuthInterceptor,
    },
    {
      provide: GRPC_INTERCEPTORS,
      multi: true,
      useClass: I18nInterceptor,
    },
    {
      provide: GRPC_INTERCEPTORS,
      multi: true,
      useClass: OrgInterceptor,
    },
    SeoService,
    RefreshService,
    GrpcService,
    AuthenticationService,
    GrpcAuthService,
    SubscriptionService,
    UploadService,
    { provide: 'windowObject', useValue: window },
  ],
  bootstrap: [AppComponent],
})

export class AppModule {
  constructor() { }
}
