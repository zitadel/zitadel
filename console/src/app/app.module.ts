import { CommonModule, registerLocaleData } from '@angular/common';
import { HttpClientModule } from '@angular/common/http';
import localeDe from '@angular/common/locales/de';
import { APP_INITIALIZER, NgModule } from '@angular/core';
import { MatNativeDateModule } from '@angular/material/core';
import { MatDialogModule } from '@angular/material/dialog';
import { MatIconModule } from '@angular/material/icon';
import { MatSelectModule } from '@angular/material/select';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { MatTooltipModule } from '@angular/material/tooltip';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { ServiceWorkerModule } from '@angular/service-worker';
import { TranslateLoader, TranslateModule } from '@ngx-translate/core';
import { AuthConfig, OAuthModule, OAuthStorage } from 'angular-oauth2-oidc';
import { QuicklinkModule } from 'ngx-quicklink';
import { from, Observable } from 'rxjs';
import { InfoOverlayModule } from 'src/app/modules/info-overlay/info-overlay.module';
import { AssetService } from 'src/app/services/asset.service';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { HasRoleModule } from './directives/has-role/has-role.module';
import { FooterModule } from './modules/footer/footer.module';
import { HeaderModule } from './modules/header/header.module';
import { KeyboardShortcutsModule } from './modules/keyboard-shortcuts/keyboard-shortcuts.module';
import { NavModule } from './modules/nav/nav.module';
import { WarnDialogModule } from './modules/warn-dialog/warn-dialog.module';
import { SignedoutComponent } from './pages/signedout/signedout.component';
import { HasRolePipeModule } from './pipes/has-role-pipe/has-role-pipe.module';
import { AdminService } from './services/admin.service';
import { AuthenticationService } from './services/authentication.service';
import { BreadcrumbService } from './services/breadcrumb.service';
import { GrpcAuthService } from './services/grpc-auth.service';
import { GrpcService } from './services/grpc.service';
import { AuthInterceptor } from './services/interceptors/auth.interceptor';
import { GRPC_INTERCEPTORS } from './services/interceptors/grpc-interceptor';
import { I18nInterceptor } from './services/interceptors/i18n.interceptor';
import { OrgInterceptor } from './services/interceptors/org.interceptor';
import { KeyboardShortcutsService } from './services/keyboard-shortcuts/keyboard-shortcuts.service';
import { ManagementService } from './services/mgmt.service';
import { NavigationService } from './services/navigation.service';
import { OverlayService } from './services/overlay/overlay.service';
import { RefreshService } from './services/refresh.service';
import { SeoService } from './services/seo.service';
import {
    StatehandlerProcessorService,
    StatehandlerProcessorServiceImpl,
} from './services/statehandler/statehandler-processor.service';
import { StatehandlerService, StatehandlerServiceImpl } from './services/statehandler/statehandler.service';
import { StorageService } from './services/storage.service';
import { ThemeService } from './services/theme.service';
import { ToastService } from './services/toast.service';

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
  declarations: [AppComponent, SignedoutComponent],
  imports: [
    AppRoutingModule,
    CommonModule,
    BrowserModule,
    HeaderModule,
    OAuthModule.forRoot({
      resourceServer: {
        allowedUrls: [
          'https://test.api.zitadel.caos.ch/caos.zitadel.auth.api.v1.AuthService',
          'https://test.api.zitadel.caos.ch/oauth/v2/userinfo',
          'https://test.api.zitadel.caos.ch/caos.zitadel.management.api.v1.ManagementService/',
          'https://preview.api.zitadel.caos.ch',
        ],
        sendAccessToken: true,
      },
    }),
    TranslateModule.forRoot({
      loader: {
        provide: TranslateLoader,
        useClass: WebpackTranslateLoader,
      },
    }),
    NavModule,
    MatNativeDateModule,
    QuicklinkModule,
    HasRoleModule,
    InfoOverlayModule,
    BrowserAnimationsModule,
    HttpClientModule,
    MatIconModule,
    MatTooltipModule,
    FooterModule,
    HasRolePipeModule,
    MatSnackBarModule,
    WarnDialogModule,
    MatSelectModule,
    MatDialogModule,
    KeyboardShortcutsModule,
    ServiceWorkerModule.register('ngsw-worker.js', { enabled: false }),
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
    OverlayService,
    SeoService,
    RefreshService,
    GrpcService,
    BreadcrumbService,
    AuthenticationService,
    GrpcAuthService,
    ManagementService,
    AdminService,
    KeyboardShortcutsService,
    AssetService,
    ToastService,
    NavigationService,
    { provide: 'windowObject', useValue: window },
  ],
  bootstrap: [AppComponent],
})
export class AppModule {
  constructor() {}
}
