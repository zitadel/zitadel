# ZITADEL Console Guide for AI Agents

## Context
The **Management Console** (`console/`) is the administrative interface for ZITADEL. It allows developers and administrators to configure organizations, projects, and users.

## Key Technology
- **Framework**: Angular (see `console/package.json` for version).
- **Language**: TypeScript.
- **State Management**: Reactive patterns with RxJS.
- **UI Component Library**: Angular Material (`@angular/material`).
- **API Communication**: gRPC-web and connectRPC via `@zitadel/client` packages.

## Architecture & Conventions

### Application Structure
```
console/src/app/
├── components/      # Shared reusable components
├── directives/      # Shared directives
├── guards/          # Route guards (auth, permissions)
├── modules/         # Feature modules (lazy-loaded)
├── pages/           # Top-level page components
├── pipes/           # Shared pipes for data transformation
├── services/        # Application services (business logic)
└── utils/           # Utility functions
```

### Service Layer Pattern
**All business logic and API calls MUST reside in services**, not components. Components should be thin orchestrators.

**Example Service Structure** (from `console/src/app/services/grpc.service.ts`):
```typescript
@Injectable({
  providedIn: 'root',  // Singleton service
})
export class GrpcService {
  // Old gRPC-web clients (legacy, maintenance mode)
  public auth!: AuthServiceClient;
  public mgmt!: ManagementServiceClient;
  public admin!: AdminServiceClient;

  // connectRPC clients from @zitadel/client/v2 (prefer for new features)
  public userNew!: ReturnType<typeof createUserServiceClient>;
  public session!: ReturnType<typeof createSessionServiceClient>;
  public featureNew!: ReturnType<typeof createFeatureServiceClient>;
  public organizationNew!: ReturnType<typeof createOrganizationServiceClient>;

  // connectRPC clients from @zitadel/client/v1 (wrap the V1 API surface via connectRPC)
  public mgmtNew!: ReturnType<typeof createManagementServiceClient>;
  public authNew!: ReturnType<typeof createAuthServiceClient>;
  public adminNew!: ReturnType<typeof createAdminServiceClient>;

  constructor(
    private readonly envService: EnvironmentService,
    private readonly authenticationService: AuthenticationService,
    private readonly translate: TranslateService,
  ) {
    // Clients are NOT initialized in the constructor.
    // Call loadAppEnvironment() before using any client.
  }
}
```

**Key Service Patterns:**
- **Dependency Injection**: Use constructor injection for all dependencies
- **RxJS Operators**: Use `switchMap`, `map`, `catchError`, `tap` for reactive streams
- **Error Handling**: Catch and transform gRPC errors into user-friendly messages
- **Interceptors**: Use for cross-cutting concerns (auth tokens, i18n headers, org context)

### State Management with RxJS
Use **reactive patterns** for state management. Avoid manual subscriptions in components when possible (use `async` pipe instead).

**Pattern:**
```typescript
// In Service
private _userData$ = new BehaviorSubject<UserData | null>(null);
public userData$ = this._userData$.asObservable();

loadUser(userId: string): Observable<UserData> {
  return this.grpcService.userNew.getUser({ id: userId }).pipe(
    tap(user => this._userData$.next(user)),
    catchError(err => {
      this.toastService.showError(err);
      return throwError(() => err);
    })
  );
}

// In Component (preferred: derive observable, consume via async pipe — no manual subscription)
export class UserComponent {
  // Drive the view from an observable; no subscribe() needed in the component.
  readonly userData$ = this.userService.loadUser('user-id');

  constructor(private readonly userService: UserService) {}
}

// In Template (preferred - automatic subscription management with no memory leaks)
<div *ngIf="userData$ | async as user">
  {{ user.name }}
</div>

// If a manual subscription is unavoidable, use takeUntilDestroyed to prevent leaks:
export class UserComponent implements OnInit {
  private readonly destroyRef = inject(DestroyRef);

  ngOnInit() {
    this.userService.loadUser('user-id').pipe(
      takeUntilDestroyed(this.destroyRef)
    ).subscribe();
  }
}
```

### gRPC-web and connectRPC Integration
- **Legacy gRPC-web clients** (`auth`, `mgmt`, `admin`): Being phased out; maintenance mode only.
- **connectRPC clients from `@zitadel/client/v2`** (`userNew`, `session`, `featureNew`, `organizationNew`): Preferred for new features targeting the V2 API surface.
- **connectRPC clients from `@zitadel/client/v1`** (`mgmtNew`, `authNew`, `adminNew`): These use the connectRPC transport but wrap the V1 API proto surface — do **not** confuse with V2 API.
- **Transport**: Uses `createGrpcWebTransport` from `@connectrpc/connect-web`
- **Interceptors**: 
  - `AuthInterceptor`: Adds bearer token to requests
  - `OrgInterceptor`: Adds organization context header
  - `I18nInterceptor`: Adds language preference header
  - `ExhaustedGrpcInterceptor`: Handles quota exhaustion errors

**Usage Example:**
```typescript
// Import from @zitadel/client
import { createUserServiceClient } from '@zitadel/client/v2';

// Use in service
this.grpcService.userNew.getUser({ id: userId }).pipe(
  map(response => this.transformUser(response)),
  catchError(this.handleError)
);
```

### Angular Material Components
**Standard Components:**
- `MatDialog` - Modal dialogs
- `MatTable` - Data tables with sorting/pagination
- `MatFormField` - Form inputs with validation
- `MatButton` - Action buttons
- `MatCard` - Content containers
- `MatSnackBar` - Toast notifications (wrapped in `ToastService`)

**Accessibility:** Always include proper labels, aria attributes, and keyboard navigation.

### Module Organization
- **Feature Modules**: Lazy-loaded modules for major features (users, projects, organizations)
- **Shared Module**: Common components, directives, pipes
- **Core Module**: Singleton services (imported once in `app.module.ts`)

### Routing & Guards
- **Route Guards**: Located in `console/src/app/guards/`
- **Auth Guard**: Ensures user is authenticated
- **Permission Guard**: Checks user permissions for routes
- **Lazy Loading**: Feature modules loaded on-demand

### Form Handling
- **Reactive Forms**: Use `FormBuilder`, `FormGroup`, `FormControl`
- **Validation**: Custom validators in `console/src/app/services/password-complexity-validator-factory.service.ts`
- **Error Messages**: Display validation errors with Material `mat-error`

### Internationalization (i18n)
- **Translation Service**: `@ngx-translate/core`
- **Language Files**: JSON translation files (see `CONTRIBUTING.md` for workflow)
- **Usage**: Inject `TranslateService` and use `translate.instant()` or `translate.get()`

## Testing Strategy

### Unit Tests
- **Status**: The `@zitadel/console` project currently has **no unit test target** configured in Nx.
- **Alternative**: Manual testing and functional UI tests cover console functionality.

### Functional UI Tests (E2E)
- **Framework**: Cypress
- **Location**: `tests/functional-ui/`
- **Run**: `pnpm nx run @zitadel/functional-ui:test`
- **Coverage**: Critical user flows (login, user management, project creation)

**When making console changes:**
1. Manually test in dev server (`pnpm nx run @zitadel/console:dev`)
2. Run functional tests for affected flows
3. Consider adding new E2E tests for new features

## Common Patterns

### Loading States
```typescript
isLoading$ = new BehaviorSubject<boolean>(false);

loadData() {
  this.isLoading$.next(true);
  this.service.getData().pipe(
    finalize(() => this.isLoading$.next(false))
  ).subscribe();
}
```

### Error Handling
```typescript
catchError((error: any) => {
  const errorMsg = this.extractErrorMessage(error);
  this.toastService.showError(errorMsg);
  return throwError(() => error);
})
```

### Pagination
```typescript
// Use MatTableDataSource with MatPaginator
dataSource = new MatTableDataSource<User>();
@ViewChild(MatPaginator) paginator!: MatPaginator;

ngAfterViewInit() {
  this.dataSource.paginator = this.paginator;
}
```

## Verified Nx Targets
- **Dev Server**: `pnpm nx run @zitadel/console:dev`
- **Build**: `pnpm nx run @zitadel/console:build`
- **Lint**: `pnpm nx run @zitadel/console:lint`
- **Generate**: `pnpm nx run @zitadel/console:generate` — runs `buf generate` to produce TypeScript/JS proto stubs in `src/app/proto/generated/`. Automatically depends on `install-proto-plugins`.
- **Install Proto Plugins**: `pnpm nx run @zitadel/console:install-proto-plugins` — downloads `protoc-gen-grpc-web` v1.5.0, `protoc-gen-js` v3.21.4, and `protoc-gen-openapiv2` v2.22.0 pre-built binaries to `.artifacts/bin/`. No Go toolchain required. Output is Nx-cached.
- **Test**: The `@zitadel/console` project currently has no `test` target configured in Nx.
- **Functional UI Tests**: Use `pnpm nx run @zitadel/functional-ui:test` (see `tests/functional-ui/AGENTS.md`).

## Best Practices
- **Prefer connectRPC V2 clients** over legacy gRPC-web V1 clients for new features
- **Always use the `async` pipe** in templates to prevent memory leaks
- **Keep components thin** - delegate to services
- **Use TypeScript strict mode** - enabled in `tsconfig.json`
- **Follow Angular style guide** - consistent naming, file structure
- **Accessibility first** - proper ARIA labels, keyboard navigation

## Cross-References
- **Functional Tests**: See `tests/functional-ui/AGENTS.md` for E2E testing
- **API Clients**: Generated from `proto/` - see `proto/AGENTS.md` and `API_DESIGN.md`
- **Contributing**: See `CONTRIBUTING.md` for i18n workflow and development setup
