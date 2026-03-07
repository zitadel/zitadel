package risk

import (
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

// EngineType determines what happens when a rule expression matches.
type EngineType string

const (
	// EngineBlock produces a blocking finding directly — no engine call.
	EngineBlock EngineType = "block"
	// EngineRateLimit checks a sliding-window counter before deciding.
	EngineRateLimit EngineType = "rate_limit"
	// EngineLLM forwards a focused context to the LLM/SLM for judgement.
	EngineLLM EngineType = "llm"
	// EngineLog emits a structured log entry (observe-mode, never blocks).
	EngineLog EngineType = "log"
)

// Rule is the YAML-level definition loaded from configuration.
type Rule struct {
	ID          string     `yaml:"id"`
	Description string     `yaml:"description"`
	Expr        string     `yaml:"expr"`
	Engine      EngineType `yaml:"engine"`

	// Finding configures the output finding when the rule matches.
	// Only Name and Block are used for the "block" engine; other engines
	// may override Block based on their own logic.
	FindingCfg RuleFinding `yaml:"finding"`

	// ContextTemplate is a Go text/template rendered with RiskContext when
	// the rule forwards to the LLM engine. Produces a compact, focused prompt.
	ContextTemplate string `yaml:"context_template"`

	// RateLimit configures the rate_limit engine.
	RateLimitCfg RuleRateLimit `yaml:"rate_limit"`
}

// RuleFinding configures the Finding emitted when a rule matches.
type RuleFinding struct {
	Name    string `yaml:"name"`
	Message string `yaml:"message"`
	Block   bool   `yaml:"block"`
}

// RuleRateLimit configures the rate_limit engine for a rule.
type RuleRateLimit struct {
	// KeyTemplate is a Go text/template rendered with RiskContext to produce
	// the counter key (e.g. "ip:{{.Current.IP}}").
	KeyTemplate string        `yaml:"key"`
	Window      time.Duration `yaml:"window"`
	Max         int           `yaml:"max"`
}

// CompiledRule is a Rule whose expression has been compiled and type-checked
// against RiskContext. It is safe to evaluate concurrently.
type CompiledRule struct {
	Rule
	program     *vm.Program
	ctxTemplate *template.Template
	keyTemplate *template.Template
}

// CompileRules compiles a slice of rule definitions against the RiskContext type.
// It returns an error if any expression fails to compile or type-check.
func CompileRules(rules []Rule) ([]CompiledRule, error) {
	compiled := make([]CompiledRule, 0, len(rules))
	for _, r := range rules {
		cr, err := compileRule(r)
		if err != nil {
			return nil, fmt.Errorf("rule %q: %w", r.ID, err)
		}
		compiled = append(compiled, cr)
	}
	return compiled, nil
}

func compileRule(r Rule) (CompiledRule, error) {
	if r.ID == "" {
		return CompiledRule{}, fmt.Errorf("rule must have an id")
	}
	if strings.TrimSpace(r.Expr) == "" {
		return CompiledRule{}, fmt.Errorf("rule must have an expression")
	}

	// Validate engine type.
	switch r.Engine {
	case EngineBlock, EngineRateLimit, EngineLLM, EngineLog:
	case "":
		return CompiledRule{}, fmt.Errorf("rule must have an engine")
	default:
		return CompiledRule{}, fmt.Errorf("unknown engine type %q", r.Engine)
	}

	// Compile expression with type-safe environment.
	program, err := expr.Compile(r.Expr,
		expr.Env(RiskContext{}),
		expr.AsBool(),
	)
	if err != nil {
		return CompiledRule{}, fmt.Errorf("compile expression: %w", err)
	}

	cr := CompiledRule{Rule: r, program: program}

	// Parse optional context template for LLM engine.
	if r.ContextTemplate != "" {
		tmpl, err := template.New(r.ID + "_ctx").Parse(r.ContextTemplate)
		if err != nil {
			return CompiledRule{}, fmt.Errorf("parse context_template: %w", err)
		}
		cr.ctxTemplate = tmpl
	}

	// Parse optional key template for rate_limit engine.
	if r.RateLimitCfg.KeyTemplate != "" {
		tmpl, err := template.New(r.ID + "_key").Parse(r.RateLimitCfg.KeyTemplate)
		if err != nil {
			return CompiledRule{}, fmt.Errorf("parse rate_limit key template: %w", err)
		}
		cr.keyTemplate = tmpl
	}

	return cr, nil
}

// Evaluate runs the compiled expression against a RiskContext and returns
// whether the rule matched.
func (cr *CompiledRule) Evaluate(rc RiskContext) (bool, error) {
	result, err := expr.Run(cr.program, rc)
	if err != nil {
		return false, fmt.Errorf("rule %q eval: %w", cr.ID, err)
	}
	matched, ok := result.(bool)
	if !ok {
		return false, fmt.Errorf("rule %q: expression returned %T, want bool", cr.ID, result)
	}
	return matched, nil
}

// RenderContextTemplate renders the LLM context template with the given RiskContext.
func (cr *CompiledRule) RenderContextTemplate(rc RiskContext) (string, error) {
	if cr.ctxTemplate == nil {
		return "", nil
	}
	var buf strings.Builder
	if err := cr.ctxTemplate.Execute(&buf, rc); err != nil {
		return "", fmt.Errorf("rule %q: render context template: %w", cr.ID, err)
	}
	return buf.String(), nil
}

// RenderKeyTemplate renders the rate limit key template with the given RiskContext.
func (cr *CompiledRule) RenderKeyTemplate(rc RiskContext) (string, error) {
	if cr.keyTemplate == nil {
		return cr.ID, nil // fallback to rule ID as key
	}
	var buf strings.Builder
	if err := cr.keyTemplate.Execute(&buf, rc); err != nil {
		return "", fmt.Errorf("rule %q: render key template: %w", cr.ID, err)
	}
	return buf.String(), nil
}
