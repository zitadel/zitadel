package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
)

// repoRoot resolves the module root the same way main() does
// (filepath.Abs(".")) when invoked from there — but `go test` runs from this
// package's own directory regardless of the caller's cwd, so this walks up
// from this file's own path instead: internal/tools/errortrace/main_test.go
// is three directories below the repo root.
func repoRoot(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	root, err := filepath.Abs(filepath.Join(filepath.Dir(thisFile), "..", "..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	return root
}

// traceGoFile runs exactly the --go-file path main() does — same
// targetsFromGoFile call, same per-target walk — but returns the result
// instead of printing it, and takes an indirections table instead of always
// using the real knownIndirections.
func traceGoFile(t *testing.T, root, goFile string, indirections map[string]funcRef) map[string]operationTrace {
	t.Helper()
	l := newLoader(root)
	targets, err := targetsFromGoFile(l, root, goFile)
	if err != nil {
		t.Fatalf("targetsFromGoFile(%s): %v", goFile, err)
	}
	if len(targets) == 0 {
		t.Fatalf("no targets found in %s", goFile)
	}
	result := make(map[string]operationTrace, len(targets))
	for _, tgt := range targets {
		w := newWalkerWithIndirections(l, root, indirections)
		w.walk(tgt.pkg, tgt.decl, []string{tgt.decl.Name.Name})
		pos := tgt.pkg.Fset.Position(tgt.decl.Name.Pos())
		result[tgt.decl.Name.Name] = operationTrace{
			Handler:    relPath(root, pos.Filename) + ":" + strconv.Itoa(pos.Line),
			Errors:     dedupeSort(w.sites),
			Unresolved: len(w.unresolved),
		}
	}
	return result
}

// TestTraceGoldenCases pins the walker's semantics against small, isolated
// fixture packages under testdata/ — one per case called out in review as
// the specific way this tool has silently gone wrong before: a direct
// throw, a throw one method call away, a call through a knownIndirections
// entry, a raw gRPC status stub (with and without a Printf verb in its
// message), the empty-ID errors.Is() sentinel that must be skipped without
// being recorded, and a dynamic call with no known target that must
// surface on Unresolved rather than vanish.
//
// The known-indirection case uses a fixture funcRef instead of the real
// knownIndirections table, so this test can't start failing just because
// internal/api/authz.CheckPermission's own internals change — it's
// pinning the lookup mechanism, not that function's business logic.
func TestTraceGoldenCases(t *testing.T) {
	root := repoRoot(t)

	cases := []struct {
		name         string
		goFile       string
		indirections map[string]funcRef
	}{
		{name: "direct_throw", goFile: "internal/tools/errortrace/testdata/direct/handler.go"},
		{name: "method_indirect", goFile: "internal/tools/errortrace/testdata/method/handler.go"},
		{
			name:   "known_indirection",
			goFile: "internal/tools/errortrace/testdata/knownindirection/handler.go",
			indirections: map[string]funcRef{
				modulePath + "/internal/tools/errortrace/testdata/knownindirection.PermCheck": {
					pkgPath: modulePath + "/internal/tools/errortrace/testdata/knownindirection",
					name:    "IndirectTarget",
				},
			},
		},
		{name: "grpc_status", goFile: "internal/tools/errortrace/testdata/grpcstatus/handler.go"},
		{name: "sentinel", goFile: "internal/tools/errortrace/testdata/sentinel/handler.go"},
		{name: "unresolved", goFile: "internal/tools/errortrace/testdata/unresolved/handler.go"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			indirections := tc.indirections
			if indirections == nil {
				indirections = map[string]funcRef{}
			}
			got := traceGoFile(t, root, filepath.Join(root, tc.goFile), indirections)
			gotJSON, err := json.MarshalIndent(got, "", "  ")
			if err != nil {
				t.Fatal(err)
			}
			gotJSON = append(gotJSON, '\n')

			goldenPath := filepath.Join("testdata", "golden", tc.name+".json")
			if os.Getenv("UPDATE_GOLDEN") == "1" {
				if err := os.WriteFile(goldenPath, gotJSON, 0o644); err != nil {
					t.Fatal(err)
				}
			}
			want, err := os.ReadFile(goldenPath)
			if err != nil {
				t.Fatalf("reading golden file %s (run `UPDATE_GOLDEN=1 go test ./internal/tools/errortrace/...` to create/refresh it): %v", goldenPath, err)
			}
			if string(gotJSON) != string(want) {
				t.Errorf("trace output for %s doesn't match golden file %s\n--- got ---\n%s\n--- want ---\n%s", tc.name, goldenPath, gotJSON, want)
			}
		})
	}
}
