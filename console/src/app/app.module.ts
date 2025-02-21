import { CommonModule, registerLocaleData } from '@angular/common';
import { HttpClientModule, HTTP_INTERCEPTORS } from '@angular/common/http';
import localeBg from '@angular/common/locales/bg';
import localeDe from '@angular/common/locales/de';
import localeCs from '@angular/common/locales/cs';
import localeEn from '@angular/common/locales/en';
import localeEs from '@angular/common/locales/es';
import localeFr from '@angular/common/locales/fr';
import localeId from '@angular/common/locales/id';
import localeIt from '@angular/common/locales/it';
import localeJa from '@angular/common/locales/ja';
import localeMk from '@angular/common/locales/mk';
import localePl from '@angular/common/locales/pl';
import localePt from '@angular/common/locales/pt';
import localeZh from '@angular/common/locales/zh';
import localeRu from '@angular/common/locales/ru';
import localeNl from '@angular/common/locales/nl';
import localeSv from '@angular/common/locales/sv';
import localeHu from '@angular/common/locales/hu';
import localeKo from '@angular/common/locales/ko';
import localeRo from '@angular/common/locales/ro';
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
import { EnvironmentService } from './services/environment.service';
import { ExhaustedService } from './services/exhausted.service';
import { GrpcAuthService } from './services/grpc-auth.service';
import { GrpcService } from './services/grpc.service';
import { AuthInterceptor } from './services/interceptors/auth.interceptor';
import { ExhaustedGrpcInterceptor } from './services/interceptors/exhausted.grpc.interceptor';
import { ExhaustedHttpInterceptor } from './services/interceptors/exhausted.http.interceptor';
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
import { LanguagesService } from './services/languages.service';
import { PosthogService } from './services/posthog.service';

registerLocaleData(localeDe);
i18nIsoCountries.registerLocale(require('i18n-iso-countries/langs/de.json'));
registerLocaleData(localeEn);
i18nIsoCountries.registerLocale(require('i18n-iso-countries/langs/en.json'));
registerLocaleData(localeEs);
i18nIsoCountries.registerLocale(require('i18n-iso-countries/langs/es.json'));
registerLocaleData(localeFr);
i18nIsoCountries.registerLocale(require('i18n-iso-countries/langs/fr.json'));
registerLocaleData(localeId);
i18nIsoCountries.registerLocale(require('i18n-iso-countries/langs/id.json'));
registerLocaleData(localeIt);
i18nIsoCountries.registerLocale(require('i18n-iso-countries/langs/it.json'));
registerLocaleData(localeJa);
i18nIsoCountries.registerLocale(require('i18n-iso-countries/langs/ja.json'));
registerLocaleData(localePl);
i18nIsoCountries.registerLocale(require('i18n-iso-countries/langs/pl.json'));
registerLocaleData(localeZh);
i18nIsoCountries.registerLocale(require('i18n-iso-countries/langs/zh.json'));
registerLocaleData(localeBg);
i18nIsoCountries.registerLocale(require('i18n-iso-countries/langs/bg.json'));
registerLocaleData(localePt);
i18nIsoCountries.registerLocale(require('i18n-iso-countries/langs/pt.json'));
registerLocaleData(localeMk);
i18nIsoCountries.registerLocale(require('i18n-iso-countries/langs/mk.json'));
registerLocaleData(localeRu);
i18nIsoCountries.registerLocale(require('i18n-iso-countries/langs/ru.json'));
registerLocaleData(localeCs);
i18nIsoCountries.registerLocale(require('i18n-iso-countries/langs/cs.json'));
registerLocaleData(localeNl);
i18nIsoCountries.registerLocale(require('i18n-iso-countries/langs/nl.json'));
registerLocaleData(localeSv);
i18nIsoCountries.registerLocale(require('i18n-iso-countries/langs/sv.json'));
registerLocaleData(localeHu);
i18nIsoCountries.registerLocale(require('i18n-iso-countries/langs/hu.json'));
registerLocaleData(localeKo);
i18nIsoCountries.registerLocale(require('i18n-iso-countries/langs/ko.json'));
registerLocaleData(localeRo);
i18nIsoCountries.registerLocale(require('i18n-iso-countries/langs/ro.json'));

export class WebpackTranslateLoader implements TranslateLoader {
  getTranslation(lang: string): Observable<any> {
    return from(import(`../assets/i18n/${lang}.json`));
  }
}

const appInitializerFn = (grpcSvc: GrpcService) => {
  return () => {
    return grpcSvc.loadAppEnvironment();
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
    EnvironmentService,
    ExhaustedService,
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
      provide: HTTP_INTERCEPTORS,
      multi: true,
      useClass: ExhaustedHttpInterceptor,
    },
    {
      provide: GRPC_INTERCEPTORS,
      multi: true,
      useClass: ExhaustedGrpcInterceptor,
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
    LanguagesService,
    PosthogService,
    { provide: 'windowObject', useValue: window },
  ],
  bootstrap: [AppComponent],
})
export class AppModule {
  constructor() {}
}
