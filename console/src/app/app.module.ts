import { CommonModule, registerLocaleData } from '@angular/common';
import { HttpClientModule } from '@angular/common/http';
import localeDe from '@angular/common/locales/de';
import localeEn from '@angular/common/locales/en';
import localeFr from '@angular/common/locales/fr';
import localeIt from '@angular/common/locales/it';
import localePl from '@angular/common/locales/pl';
import localeZh from '@angular/common/locales/zh';
import { APP_INITIALIZER, NgModule } from '@angular/core';
import { MatNativeDateModule } from '@angular/material/core';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyDialogModule as MatDialogModule } from '@angular/material/legacy-dialog';
import { MatLegacySelectModule as MatSelectModule } from '@angular/material/legacy-select';
import { MatLegacySnackBarModule as MatSnackBarModule } from '@angular/material/legacy-snack-bar';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { ServiceWorkerModule } from '@angular/service-worker';
import { TranslateLoader, TranslateModule } from '@ngx-translate/core';
import { AuthConfig, OAuthModule, OAuthStorage } from 'angular-oauth2-oidc';
import * as i18nIsoCountries from 'i18n-iso-countries';
import { from, Observable } from 'rxjs';
import { AuthGuard } from 'src/app/guards/auth.guard';
import { RoleGuard } from 'src/app/guards/role.guard';
import { UserGuard } from 'src/app/guards/user.guard';
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
i18nIsoCountries.registerLocale(require('i18n-iso-countries/langs/de.json'));
registerLocaleData(localeZh);
i18nIsoCountries.registerLocale(require('i18n-iso-countries/langs/zh.json'));
registerLocaleData(localeFr);
i18nIsoCountries.registerLocale(require('i18n-iso-countries/langs/fr.json'));
registerLocaleData(localeIt);
i18nIsoCountries.registerLocale(require('i18n-iso-countries/langs/it.json'));
registerLocaleData(localePl);
i18nIsoCountries.registerLocale(require('i18n-iso-countries/langs/pl.json'));
registerLocaleData(localeEn);
i18nIsoCountries.registerLocale(require('i18n-iso-countries/langs/en.json'));

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
  requireHttps: false,
};

@NgModule({
  declarations: [AppComponent],
  imports: [
    AppRoutingModule,
    CommonModule,
    BrowserModule,
    HeaderModule,
    OAuthModule.forRoot(),
    TranslateModule.forRoot({
      loader: {
        provide: TranslateLoader,
        useClass: WebpackTranslateLoader,
      },
    }),
    NavModule,
    MatNativeDateModule,
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
    AuthGuard,
    RoleGuard,
    UserGuard,
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
