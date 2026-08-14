// Command errortrace computes, for one or more RPC handler methods, the set
// of zerrors.Throw*() call sites reachable via the real Go call graph. It
// replaces the hand-tracing step described in
// apps/docs/scripts/generate-endpoint-errors.ts: instead of a human reading
// code and writing down {id, file, line} triples, this walks the type-checked
// AST starting at a handler method and follows every statically resolvable
// call (direct function calls, and method calls on concrete, non-interface
// receivers) until it hits a zerrors.Throw* call or runs out of module code
// to follow.
//
// Two entry points, matching how a dev actually thinks about "what do I want
// traced":
//
//	go run ./internal/tools/errortrace --go-file internal/api/grpc/user/v2/metadata.go
//	go run ./internal/tools/errortrace --proto proto/zitadel/user/v2/user_service.proto
//
// Output (stdout) is JSON in the exact shape already used by the hand-written
// files under apps/docs/scripts/endpoint-error-tracing/<service>/*.json, so
// nothing downstream (generate-endpoint-errors.ts) needs to change.
//
// Calls made through a function-typed variable/field/parameter (e.g. the
// PermissionCheck closures) can't be resolved to one concrete target by
// static analysis alone. Rather than guess, such calls are checked against a
// small table of known chokepoints (see knownIndirections below) and
// otherwise reported as unresolved on stderr, never silently dropped.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/tools/go/packages"
)

const modulePath = "github.com/zitadel/zitadel"

// throwFuncs is every zerrors.Throw<Kind>[f] function. All 24 share the
// convention that the error ID literal is their 2nd positional argument
// (confirmed against internal/zerrors/*.go: Throw<Kind>(parent, id, message)
// and Throw<Kind>f(parent, id, format, a...) both put id at index 1).
var throwFuncs = func() map[string]bool {
	kinds := []string{
		"InvalidArgument", "NotFound", "AlreadyExists", "PermissionDenied",
		"Internal", "PreconditionFailed", "Unauthenticated", "Unimplemented",
		"DeadlineExceeded", "Unavailable", "ResourceExhausted", "Unknown",
	}
	m := make(map[string]bool, len(kinds)*2+1)
	for _, k := range kinds {
		m["Throw"+k] = true
		m["Throw"+k+"f"] = true
	}
	// ThrowError doesn't fit the Throw<Kind> naming pattern (it's not "Throw"
	// + a Kind name) but has the exact same (parent, id, message) shape and
	// always maps to KindUnknown (internal/zerrors/zerror.go). No ThrowErrorf
	// variant exists.
	m["ThrowError"] = true
	return m
}()

// funcRef names a free function to jump into instead of a call site we can't
// resolve locally.
type funcRef struct {
	pkgPath string
	name    string
}

// knownIndirections maps the fully-qualified name of a named function type
// to the real function that ends up handling any call through a value of
// that type. Seeded from chains already verified by hand while tracing the
// 11 operations under apps/docs/scripts/endpoint-error-tracing/user-v2/: every
// PermissionCheck closure, regardless of which factory built it
// (NewPermissionCheckUserWrite, checkPermissionOnUser, ...), bottoms out at
// authz.CheckPermission.
var knownIndirections = map[string]funcRef{
	"github.com/zitadel/zitadel/internal/command.PermissionCheck": {
		pkgPath: "github.com/zitadel/zitadel/internal/api/authz", name: "CheckPermission",
	},
	"github.com/zitadel/zitadel/internal/domain.PermissionCheck": {
		pkgPath: "github.com/zitadel/zitadel/internal/api/authz", name: "CheckPermission",
	},
	// The three concrete CheckPermissionOrganization{Write,Create,Delete}
	// methods (internal/command/permission_checks.go) all just call
	// newPermissionCheck with a different permission string — same downstream
	// errors regardless of which one a given operation actually passes.
	"github.com/zitadel/zitadel/internal/command.OrganizationPermissionCheck": {
		pkgPath: "github.com/zitadel/zitadel/internal/api/authz", name: "CheckPermission",
	},
}

type errorSite struct {
	ID        string `json:"id"`
	File      string `json:"file"`
	Line      int    `json:"line"`
	Reasoning string `json:"reasoning"`
	// Message is only set for synthetic sites (id starts with "GRPC-") that
	// have no entry in the error catalog to pull a message from — a raw
	// status.Errorf(codes.X, ...) call, unlike a zerrors.Throw* call, was
	// never scanned into that catalog in the first place, so the message
	// has to travel with the site itself instead of being looked up later.
	Message string `json:"message,omitempty"`
}

type operationTrace struct {
	Handler string      `json:"handler"`
	Errors  []errorSite `json:"errors"`
}

func main() {
	goFile := flag.String("go-file", "", "trace every exported handler method defined in this Go file")
	protoFile := flag.String("proto", "", "trace every RPC declared in this service .proto file")
	goPackage := flag.String("go-package", "", "Go package (relative to repo root) holding the proto service's handlers; defaults to the proto path convention proto/zitadel/<x>/v2/... -> internal/api/grpc/<x>/v2")
	repoRoot := flag.String("repo-root", ".", "repository root (module root); defaults to the current directory")
	out := flag.String("out", "", "write JSON here instead of stdout")
	flag.Parse()

	if (*goFile == "") == (*protoFile == "") {
		fmt.Fprintln(os.Stderr, "errortrace: exactly one of --go-file or --proto is required")
		os.Exit(2)
	}

	root, err := filepath.Abs(*repoRoot)
	if err != nil {
		fmt.Fprintln(os.Stderr, "errortrace:", err)
		os.Exit(1)
	}

	l := newLoader(root)

	var targets []target
	switch {
	case *goFile != "":
		targets, err = targetsFromGoFile(l, root, *goFile)
	case *protoFile != "":
		targets, err = targetsFromProto(l, root, *protoFile, *goPackage)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "errortrace:", err)
		os.Exit(1)
	}
	if len(targets) == 0 {
		fmt.Fprintln(os.Stderr, "errortrace: no operations found")
		os.Exit(1)
	}

	result := make(map[string]operationTrace, len(targets))
	unresolvedTotal := 0
	for _, t := range targets {
		w := newWalker(l, root)
		w.walk(t.pkg, t.decl, []string{t.decl.Name.Name})
		sites := dedupeSort(w.sites)
		pos := t.pkg.Fset.Position(t.decl.Name.Pos())
		result[t.decl.Name.Name] = operationTrace{
			Handler: relPath(root, pos.Filename) + ":" + strconv.Itoa(pos.Line),
			Errors:  sites,
		}
		for _, u := range w.unresolved {
			fmt.Fprintf(os.Stderr, "errortrace: %s: unresolved dynamic call to a value of type %s (not in the known-indirection table)\n", u.pos, u.typeName)
		}
		unresolvedTotal += len(w.unresolved)
		fmt.Fprintf(os.Stderr, "errortrace: %s: %d error site(s) found\n", t.decl.Name.Name, len(sites))
	}
	if unresolvedTotal > 0 {
		fmt.Fprintf(os.Stderr, "errortrace: %d unresolved dynamic call site(s) total — extend knownIndirections in main.go if these matter\n", unresolvedTotal)
	}

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "errortrace:", err)
		os.Exit(1)
	}
	data = append(data, '\n')

	if *out == "" {
		os.Stdout.Write(data)
		return
	}
	if err := os.WriteFile(*out, data, 0o644); err != nil {
		fmt.Fprintln(os.Stderr, "errortrace:", err)
		os.Exit(1)
	}
}

// --- entry point resolution -------------------------------------------------

type target struct {
	pkg  *packages.Package
	decl *ast.FuncDecl
}

// targetsFromGoFile traces every exported method with a receiver that is
// defined directly in the given file.
func targetsFromGoFile(l *loader, root, goFile string) ([]target, error) {
	abs := goFile
	if !filepath.IsAbs(abs) {
		abs = filepath.Join(root, goFile)
	}
	pkg, err := l.loadFile(abs)
	if err != nil {
		return nil, err
	}
	var targets []target
	for _, f := range pkg.Syntax {
		pos := pkg.Fset.Position(f.Pos())
		if pos.Filename != abs {
			continue
		}
		for _, decl := range f.Decls {
			fd, ok := decl.(*ast.FuncDecl)
			if !ok || fd.Recv == nil || !fd.Name.IsExported() {
				continue
			}
			targets = append(targets, target{pkg: pkg, decl: fd})
		}
	}
	return targets, nil
}

var rpcNameRe = regexp.MustCompile(`\brpc\s+(\w+)\s*\(`)

// targetsFromProto parses `rpc <Name>(` declarations out of the given
// service .proto file (the same convention already used by
// apps/docs/scripts/generate-endpoint-errors.ts's declared-response parser)
// and resolves each name to its Go handler method.
func targetsFromProto(l *loader, root, protoFile, goPackage string) ([]target, error) {
	protoAbs := protoFile
	if !filepath.IsAbs(protoAbs) {
		protoAbs = filepath.Join(root, protoFile)
	}
	text, err := os.ReadFile(protoAbs)
	if err != nil {
		return nil, err
	}
	names := map[string]bool{}
	for _, m := range rpcNameRe.FindAllStringSubmatch(string(text), -1) {
		names[m[1]] = true
	}
	if len(names) == 0 {
		return nil, fmt.Errorf("%s: no `rpc Name(...)` declarations found — is this a service file? (message-only files like metadata.proto have none)", protoFile)
	}

	pkgDir := goPackage
	if pkgDir == "" {
		pkgDir, err = derivePackageDir(root, protoAbs)
		if err != nil {
			return nil, err
		}
	}
	pkg, err := l.loadDir(filepath.Join(root, pkgDir))
	if err != nil {
		return nil, err
	}

	byName := map[string]*ast.FuncDecl{}
	for _, f := range pkg.Syntax {
		for _, decl := range f.Decls {
			if fd, ok := decl.(*ast.FuncDecl); ok && fd.Recv != nil {
				byName[fd.Name.Name] = fd
			}
		}
	}

	var targets []target
	var missing []string
	for name := range names {
		fd, ok := byName[name]
		if !ok {
			missing = append(missing, name)
			continue
		}
		targets = append(targets, target{pkg: pkg, decl: fd})
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		fmt.Fprintf(os.Stderr, "errortrace: no handler method found in %s for: %s\n", pkgDir, strings.Join(missing, ", "))
	}
	sort.Slice(targets, func(i, j int) bool { return targets[i].decl.Name.Name < targets[j].decl.Name.Name })
	return targets, nil
}

// derivePackageDir maps proto/zitadel/<x>/v2/<name>_service.proto to
// internal/api/grpc/<x>/v2, the convention confirmed across user, org and
// session v2 services.
func derivePackageDir(root, protoAbs string) (string, error) {
	rel, err := filepath.Rel(filepath.Join(root, "proto", "zitadel"), protoAbs)
	if err != nil {
		return "", err
	}
	parts := strings.Split(filepath.ToSlash(rel), "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("%s: doesn't match proto/zitadel/<service>/<version>/... — pass --go-package explicitly", protoAbs)
	}
	return filepath.Join("internal", "api", "grpc", parts[0], parts[1]), nil
}

// --- package loading ---------------------------------------------------

// loader lazily loads and caches packages by import path, sharing one
// token.FileSet so positions from different loads stay comparable. Loading
// is on-demand rather than a single NeedDeps load of the whole module: the
// walker only ever asks for packages it actually recurses into (it never
// follows dynamic dispatch outside knownIndirections), so in practice only a
// handful of packages get loaded per trace instead of the whole transitive
// dependency graph.
type loader struct {
	fset *token.FileSet
	cfg  *packages.Config
	byID map[string]*packages.Package // import path -> package
}

func newLoader(root string) *loader {
	fset := token.NewFileSet()
	return &loader{
		fset: fset,
		cfg: &packages.Config{
			Mode: packages.NeedName | packages.NeedFiles | packages.NeedImports |
				packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax,
			Dir:  root,
			Fset: fset,
		},
		byID: map[string]*packages.Package{},
	}
}

func (l *loader) loadDir(dir string) (*packages.Package, error) {
	pkgs, err := packages.Load(l.cfg, "./"+relPath(l.cfg.Dir, dir))
	if err != nil {
		return nil, err
	}
	return l.first(dir, pkgs)
}

func (l *loader) loadFile(absPath string) (*packages.Package, error) {
	pkgs, err := packages.Load(l.cfg, "file="+absPath)
	if err != nil {
		return nil, err
	}
	return l.first(absPath, pkgs)
}

func (l *loader) loadImportPath(importPath string) (*packages.Package, error) {
	if pkg, ok := l.byID[importPath]; ok {
		return pkg, nil
	}
	pkgs, err := packages.Load(l.cfg, importPath)
	if err != nil {
		return nil, err
	}
	pkg, err := l.first(importPath, pkgs)
	if err != nil {
		return nil, err
	}
	l.byID[importPath] = pkg
	return pkg, nil
}

func (l *loader) first(what string, pkgs []*packages.Package) (*packages.Package, error) {
	if len(pkgs) == 0 {
		return nil, fmt.Errorf("no package found for %s", what)
	}
	pkg := pkgs[0]
	if len(pkg.Errors) > 0 {
		return nil, fmt.Errorf("%s: %v", what, pkg.Errors[0])
	}
	l.byID[pkg.PkgPath] = pkg
	return pkg, nil
}

// findFuncDecl looks up a function/method by name (and receiver type name,
// "" for free functions) within an already-loaded package. Matching by name
// rather than *types.Object identity sidesteps a real pitfall: objects
// resolved from one package's TypesInfo aren't guaranteed to be pointer-equal
// to objects from a separately loaded package/module (export-data-based
// imports don't carry the same identity as a fresh source load), so pointer
// comparison across the lazily-loaded package boundary would be unreliable.
func findFuncDecl(pkg *packages.Package, name, recv string) *ast.FuncDecl {
	for _, f := range pkg.Syntax {
		for _, decl := range f.Decls {
			fd, ok := decl.(*ast.FuncDecl)
			if !ok || fd.Name.Name != name {
				continue
			}
			if receiverTypeName(fd) == recv {
				return fd
			}
		}
	}
	return nil
}

func receiverTypeName(fd *ast.FuncDecl) string {
	if fd.Recv == nil || len(fd.Recv.List) == 0 {
		return ""
	}
	expr := fd.Recv.List[0].Type
	if star, ok := expr.(*ast.StarExpr); ok {
		expr = star.X
	}
	if id, ok := expr.(*ast.Ident); ok {
		return id.Name
	}
	return ""
}

// --- the walk ------------------------------------------------------------

type unresolvedCall struct {
	pos      string
	typeName string
}

type walker struct {
	l          *loader
	root       string
	visited    map[string]bool // "pkgPath#name#recv" cycle/dedup guard, scoped to one operation's trace
	sites      []errorSite
	unresolved []unresolvedCall
}

func newWalker(l *loader, root string) *walker {
	return &walker{l: l, root: root, visited: map[string]bool{}}
}

func funcKey(pkgPath, name, recv string) string { return pkgPath + "#" + name + "#" + recv }

func (w *walker) walk(pkg *packages.Package, decl *ast.FuncDecl, trail []string) {
	if decl.Body == nil {
		return
	}
	w.walkBody(pkg, decl.Body, funcKey(pkg.PkgPath, decl.Name.Name, receiverTypeName(decl)), trail)
}

func (w *walker) walkBody(pkg *packages.Package, body *ast.BlockStmt, visitKey string, trail []string) {
	if w.visited[visitKey] {
		return
	}
	w.visited[visitKey] = true

	ast.Inspect(body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		w.handleCall(pkg, call, trail)
		return true
	})
}

// findPackageLevelFuncLit looks for `var Name = func(...) {...}` (or `var (
// Name = func...)` inside a grouped block) at package scope — the pattern
// used for e.g. internal/command/user_human_password.go's ErrPasswordInvalid.
// A *types.Var resolved to this has no FuncDecl to find via the normal path,
// since it's a literal assigned to a var, not a declared function.
func findPackageLevelFuncLit(pkg *packages.Package, name string) *ast.FuncLit {
	for _, f := range pkg.Syntax {
		for _, decl := range f.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.VAR {
				continue
			}
			for _, spec := range gd.Specs {
				vs, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}
				for i, n := range vs.Names {
					if n.Name != name || i >= len(vs.Values) {
						continue
					}
					if lit, ok := vs.Values[i].(*ast.FuncLit); ok {
						return lit
					}
				}
			}
		}
	}
	return nil
}

func (w *walker) handleCall(pkg *packages.Package, call *ast.CallExpr, trail []string) {
	switch fun := call.Fun.(type) {
	case *ast.SelectorExpr:
		if sel, ok := pkg.TypesInfo.Selections[fun]; ok {
			// Selections cover both method calls (Obj() is a *types.Func) and
			// field accesses (Obj() is a *types.Var) — e.g. c.checkPermission(...)
			// where checkPermission is a struct field holding a function value,
			// not a declared method. That's dynamic dispatch exactly like a
			// bare identifier of function type, so it must go through the same
			// handleDynamic path (and knownIndirections table) rather than
			// being silently dropped here.
			if fnObj, ok := sel.Obj().(*types.Func); ok {
				if isInterfaceRecv(sel.Recv()) {
					w.handleDynamic(pkg, sel.Recv(), call)
					return
				}
				w.followFunc(pkg, fnObj, call, trail)
				return
			}
			if fieldObj, ok := sel.Obj().(*types.Var); ok {
				w.handleDynamic(pkg, fieldObj.Type(), call)
			}
			return
		}
		// Package-qualified call, e.g. zerrors.ThrowInvalidArgument(...) or
		// domain.RenderConfirmURLTemplate(...).
		if use, ok := pkg.TypesInfo.Uses[fun.Sel]; ok {
			if fnObj, ok := use.(*types.Func); ok {
				w.followFunc(pkg, fnObj, call, trail)
			}
		}
	case *ast.Ident:
		use, ok := pkg.TypesInfo.Uses[fun]
		if !ok {
			return
		}
		switch obj := use.(type) {
		case *types.Func:
			w.followFunc(pkg, obj, call, trail)
		case *types.Var:
			if lit, litPkg := w.findVarFuncLit(obj); lit != nil {
				key := funcKey(obj.Pkg().Path(), obj.Name(), "")
				w.walkBody(litPkg, lit.Body, key, append(trail, obj.Name()))
				return
			}
			w.handleDynamic(pkg, obj.Type(), call)
		}
	case *ast.CallExpr:
		// An immediately-invoked result, e.g.
		// c.NewPermissionCheckUserWrite(ctx, false)(resourceOwner, userID) —
		// the outer call's Fun is itself a call. We can't know which closure
		// comes back without evaluating fun, so treat it like any other
		// dynamic dispatch keyed on the static (named) return type.
		if t := pkg.TypesInfo.TypeOf(fun); t != nil {
			w.handleDynamic(pkg, t, call)
		}
	}
}

// followFunc handles a statically resolved callee: record it if it's a
// zerrors.Throw* call, otherwise load its defining package (if it's part of
// this module) and recurse into its body.
func (w *walker) followFunc(pkg *packages.Package, fnObj *types.Func, call *ast.CallExpr, trail []string) {
	if fnObj.Pkg() == nil {
		return // builtin
	}
	pkgPath := fnObj.Pkg().Path()
	if pkgPath == "github.com/zitadel/zitadel/internal/zerrors" && throwFuncs[fnObj.Name()] {
		w.recordThrow(pkg, call, trail)
		return
	}
	// The gRPC-native equivalent of a Throw call: status.Errorf(codes.X, ...)
	// / status.Error(codes.X, ...). Common in stub handlers
	// (status.Errorf(codes.Unimplemented, "method X not implemented")) that
	// never touch zerrors at all — a real, deterministic error with no
	// zerrors ID to look up, so it's recorded as its own synthetic site
	// rather than silently having nothing to show.
	if pkgPath == "google.golang.org/grpc/status" && (fnObj.Name() == "Errorf" || fnObj.Name() == "Error") {
		w.recordRawStatus(pkg, call, trail)
		return
	}
	if !strings.HasPrefix(pkgPath, modulePath) {
		return // stdlib / third-party: no source to follow, and never a Throw site itself
	}
	recv := ""
	if sig, ok := fnObj.Type().(*types.Signature); ok && sig.Recv() != nil {
		recv = bareTypeName(sig.Recv().Type())
	}
	w.recurseInto(pkgPath, fnObj.Name(), recv, append(trail, fnObj.Name()))
}

// handleDynamic is reached for a call through an interface method or a
// function-typed variable/field/parameter — something static analysis can't
// resolve to one concrete target. It checks knownIndirections before giving
// up.
func (w *walker) handleDynamic(pkg *packages.Package, t types.Type, call *ast.CallExpr) {
	name := namedTypeName(t)
	if ref, ok := knownIndirections[name]; ok {
		w.recurseInto(ref.pkgPath, ref.name, "", []string{ref.name + " (via " + name + ")"})
		return
	}
	// Calls through a stdlib/third-party interface (context.Context,
	// telemetry spans, ...) are extremely common (every traced function
	// carries tracing boilerplate) and can never themselves reach a
	// zerrors.Throw*, since that package isn't part of this module. Only
	// calls through an in-module type are worth surfacing as a real gap.
	if name == "" || !strings.HasPrefix(name, modulePath) {
		return
	}
	pos := pkg.Fset.Position(call.Pos())
	w.unresolved = append(w.unresolved, unresolvedCall{
		pos:      relPath(w.root, pos.Filename) + ":" + strconv.Itoa(pos.Line),
		typeName: name,
	})
}

// findVarFuncLit resolves a *types.Var to a package-level func-literal
// initializer, if it's one of those rather than a local/parameter var. It
// returns the FuncLit's package (loaded fresh if needed, same as
// recurseInto) so the caller has the right TypesInfo to walk its body with.
func (w *walker) findVarFuncLit(obj *types.Var) (*ast.FuncLit, *packages.Package) {
	if obj.Pkg() == nil || !strings.HasPrefix(obj.Pkg().Path(), modulePath) {
		return nil, nil
	}
	varPkg, err := w.l.loadImportPath(obj.Pkg().Path())
	if err != nil {
		return nil, nil
	}
	lit := findPackageLevelFuncLit(varPkg, obj.Name())
	if lit == nil {
		return nil, nil
	}
	return lit, varPkg
}

func (w *walker) recurseInto(pkgPath, name, recv string, trail []string) {
	pkg, err := w.l.loadImportPath(pkgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "errortrace: could not load %s to follow %s: %v\n", pkgPath, name, err)
		return
	}
	decl := findFuncDecl(pkg, name, recv)
	if decl == nil {
		fmt.Fprintf(os.Stderr, "errortrace: %s.%s (recv=%q) not found in loaded source\n", pkgPath, name, recv)
		return
	}
	w.walk(pkg, decl, trail)
}

func (w *walker) recordThrow(pkg *packages.Package, call *ast.CallExpr, trail []string) {
	if len(call.Args) < 2 {
		return
	}
	lit, ok := call.Args[1].(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		pos := pkg.Fset.Position(call.Pos())
		w.unresolved = append(w.unresolved, unresolvedCall{
			pos:      relPath(w.root, pos.Filename) + ":" + strconv.Itoa(pos.Line),
			typeName: "zerrors call with a non-literal id argument",
		})
		return
	}
	id, err := strconv.Unquote(lit.Value)
	if err != nil {
		return
	}
	pos := pkg.Fset.Position(call.Pos())
	if id == "" {
		// A real Throw site always has a non-empty ID — that's the whole
		// point of the ID (an addressable, greppable tag). An empty one
		// (e.g. zerrors.ThrowNotFound(nil, "", "")) is the errors.Is()
		// sentinel-value idiom: the constructed error is only used for a
		// type comparison, never actually returned to a caller. See
		// internal/command/org_domain.go:44 for the pattern this guards
		// against.
		fmt.Fprintf(os.Stderr, "errortrace: %s: skipped a zerrors.Throw call with an empty ID (looks like an errors.Is() sentinel, not a real returned error)\n",
			relPath(w.root, pos.Filename)+":"+strconv.Itoa(pos.Line))
		return
	}
	w.sites = append(w.sites, errorSite{
		ID:        id,
		File:      relPath(w.root, pos.Filename),
		Line:      pos.Line,
		Reasoning: strings.Join(trail, " -> "),
	})
}

// recordRawStatus handles status.Errorf(codes.X, format, a...) /
// status.Error(codes.X, msg) — the first argument must resolve to a real
// constant in google.golang.org/grpc/codes (not just any identifier named
// like a code), so an unrelated call that happens to pass an argument named
// "codes.Foo" from some other package can't be mistaken for this. There's
// no zerrors ID here, so a synthetic one ("GRPC-<CODE NAME>") is used
// instead — clearly distinguishable from a real zerrors ID by its shape,
// and matched specially in generate-endpoint-errors.ts rather than looked
// up in the error catalog, since it was never scanned into that catalog.
// grpcCodeToStatusKey maps a google.golang.org/grpc/codes constant name to
// the exact key GRPC_STATUS (apps/docs/lib/grpc-status.ts) uses for it.
// Deliberately not derived by upper-casing the Go name: several of these
// don't round-trip (NotFound -> NOT_FOUND needs an inserted underscore;
// Go spells Canceled with one L, GRPC_STATUS's key has two). Getting this
// wrong doesn't just mean the code's status text is missing — combined with
// the wasTraced filter in generate-endpoint-errors.ts, a genuinely-reachable
// site that fails this lookup gets counted as "checked, confirmed empty"
// instead of surfacing as unmatched, which is a worse failure than a gap.
var grpcCodeToStatusKey = map[string]string{
	"OK":                 "OK",
	"Canceled":           "CANCELLED",
	"Unknown":            "UNKNOWN",
	"InvalidArgument":    "INVALID_ARGUMENT",
	"DeadlineExceeded":   "DEADLINE_EXCEEDED",
	"NotFound":           "NOT_FOUND",
	"AlreadyExists":      "ALREADY_EXISTS",
	"PermissionDenied":   "PERMISSION_DENIED",
	"ResourceExhausted":  "RESOURCE_EXHAUSTED",
	"FailedPrecondition": "FAILED_PRECONDITION",
	"Aborted":            "ABORTED",
	"OutOfRange":         "OUT_OF_RANGE",
	"Unimplemented":      "UNIMPLEMENTED",
	"Internal":           "INTERNAL",
	"Unavailable":        "UNAVAILABLE",
	"DataLoss":           "DATA_LOSS",
	"Unauthenticated":    "UNAUTHENTICATED",
}

func (w *walker) recordRawStatus(pkg *packages.Package, call *ast.CallExpr, trail []string) {
	if len(call.Args) < 1 {
		return
	}
	sel, ok := call.Args[0].(*ast.SelectorExpr)
	if !ok {
		return
	}
	use, ok := pkg.TypesInfo.Uses[sel.Sel]
	if !ok {
		return
	}
	constObj, ok := use.(*types.Const)
	if !ok || constObj.Pkg() == nil || constObj.Pkg().Path() != "google.golang.org/grpc/codes" {
		return
	}

	message := ""
	if len(call.Args) >= 2 {
		if lit, ok := call.Args[1].(*ast.BasicLit); ok && lit.Kind == token.STRING {
			if s, err := strconv.Unquote(lit.Value); err == nil {
				message = s
			}
		}
	}

	pos := pkg.Fset.Position(call.Pos())
	statusKey, ok := grpcCodeToStatusKey[sel.Sel.Name]
	if !ok {
		// codes.<Name> resolved to a real constant in the codes package (the
		// checks above confirm that), but isn't one of the 17 known ones —
		// shouldn't happen in practice, but fail loud into unresolved rather
		// than emit an ID nothing can look up.
		w.unresolved = append(w.unresolved, unresolvedCall{
			pos:      relPath(w.root, pos.Filename) + ":" + strconv.Itoa(pos.Line),
			typeName: "unrecognized grpc/codes constant: " + sel.Sel.Name,
		})
		return
	}
	w.sites = append(w.sites, errorSite{
		ID:        "GRPC-" + statusKey,
		File:      relPath(w.root, pos.Filename),
		Line:      pos.Line,
		Reasoning: strings.Join(trail, " -> "),
		Message:   message,
	})
}

func isInterfaceRecv(t types.Type) bool {
	_, ok := t.Underlying().(*types.Interface)
	return ok
}

// namedTypeName returns "<import path>.<name>" for a named type (following
// through a leading pointer), or "" if t isn't a named type — e.g. a raw
// func literal type has no stable name to key knownIndirections on.
func namedTypeName(t types.Type) string {
	if ptr, ok := t.(*types.Pointer); ok {
		t = ptr.Elem()
	}
	named, ok := t.(*types.Named)
	if !ok || named.Obj().Pkg() == nil {
		return ""
	}
	return named.Obj().Pkg().Path() + "." + named.Obj().Name()
}

// bareTypeName returns just the type's own name (no package qualification),
// matching how findFuncDecl reads a receiver straight off the AST (an
// unqualified identifier, since within its own package a type is always
// referred to unqualified).
func bareTypeName(t types.Type) string {
	if ptr, ok := t.(*types.Pointer); ok {
		t = ptr.Elem()
	}
	named, ok := t.(*types.Named)
	if !ok {
		return ""
	}
	return named.Obj().Name()
}

func relPath(root, path string) string {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return path
	}
	return filepath.ToSlash(rel)
}

func dedupeSort(sites []errorSite) []errorSite {
	seen := map[string]bool{}
	out := make([]errorSite, 0, len(sites))
	for _, s := range sites {
		key := s.ID + "#" + s.File + "#" + strconv.Itoa(s.Line)
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, s)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].File != out[j].File {
			return out[i].File < out[j].File
		}
		return out[i].Line < out[j].Line
	})
	return out
}
