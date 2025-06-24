# When AI Meets Mockery: How I Stopped Fighting Mock Generation

You know that moment when you ask an AI coding assistant to generate some mocks for your Go tests, and it confidently produces what looks like reasonable code? You run the tests and discover half of them are flaky, the mocks don't follow your project patterns, and you spend more time fixing the AI's work than writing the mocks yourself.

I hit this wall repeatedly. Instead of trying to teach AI to write better mocks (spoiler: that approach was a dead end), I flipped the problem entirely. What if I stopped asking AI to generate code and started giving it better tools to orchestrate instead?

## The Core Insight

Here's what changed everything: AI is surprisingly good at understanding workflows and coordinating tools. It's terrible at generating consistent, framework-compliant code from scratch.

This led me to experiment with three different approaches to reliable mock generation, all built around [Mockery](https://github.com/vektra/mockery) as the underlying engine. I built working implementations of each approach, tested them with real interfaces, and documented what actually worked versus what sounded good in theory.

The results were illuminating, and honestly, not what I expected.

## The Three Experiments

### Experiment 1: Prompt Engineering (The Optimistic Phase)

**Location**: [`examples/prompt-only/`](examples/prompt-only/)

**The Theory**: Maybe I just needed to write really, really detailed prompts.

I crafted exhaustive guidelines covering mock patterns, naming conventions, and framework integration. The AI would read these instructions and generate mocks accordingly, essentially turning it into an expensive template engine.

**What Actually Happened**:

- Worked fine for simple cases
- Zero infrastructure overhead
- Broke down completely with complex interface hierarchies
- Quality varied wildly based on the AI's interpretation
- Debugging meant rewriting prompts until something worked

**Reality Check**: This is the "training wheels" approach. Gets you started but doesn't scale past toy examples.

### Experiment 2: Direct Tooling (The Sweet Spot)

**Locations**: [`examples/command-line-only/`](examples/command-line-only/) and [`examples/makefile/`](examples/makefile/)

**The Theory**: Let AI orchestrate Mockery commands instead of generating mock code.

I built this two ways. The command-line version has AI run complex mockery commands directly:

```bash
mockery --name=UserRepository --dir=./internal/interfaces --output=./mocks
```

The Makefile version simplifies this to:

```bash
make generate-mocks
```

**What Actually Happened**:

- Both approaches just work
- Leverage proven Mockery tooling with full feature access
- Type-safe, compiler-verified output
- Fast execution with minimal overhead
- AI focuses on orchestration, not code generation

**The Makefile Advantage**: After building both, the Makefile version won by a clear margin:

- Simpler AI prompts: "run `make generate-mocks`" vs. complex mockery syntax
- Self-documenting: `make help` shows all options
- Composable: Easy to add `make test`, `make clean-mocks`, etc.
- Consistent: Same commands across all projects
- **Massive token savings**: 82% reduction in AI conversation costs compared to the Command-line approach.

**Real-World Token Impact**:
For a typical AI conversation generating mocks:

- **Command-line approach**: ~2,330 tokens (2,130 context + 200 response)
- **Makefile approach**: ~420 tokens (370 context + 50 response)
- **Result**: 82% token reduction per conversation

This adds up fast when you're generating mocks regularly.

**The Catch**: Requires AI agents that can execute shell commands. Not all of them can, which limits where you can use this approach.

### Experiment 3: Structured Tooling (The Over-Engineering Special)

**Location**: [`mcp-server/`](mcp-server/)

**The Theory**: Build a production-grade service that wraps Mockery with proper APIs, validation, and all the enterprise features.

I implemented a Go-based MCP server with AST parsing, interface discovery, Docker deployment, comprehensive error handling, and enough configuration options to make your head spin.

**What Actually Happened**:

- It works beautifully
- Production-ready architecture with proper isolation
- Supports multiple AI agents simultaneously
- Advanced features like dependency analysis
- Took significantly longer to build than the other approaches

**The Catch**: This solves problems most teams don't have yet. Unless you're running multiple AI agents across large codebases, it's impressive overkill.

**Reality Check**: I built this because I could, not because I needed to. It's the kind of solution that looks great in architecture reviews but makes you question your life choices during 2 AM debugging sessions.

## The Surprising Winner

After implementing all approaches with real Go interfaces ([`UserRepository`](examples/command-line-only/internal/interfaces/repository.go), [`EmailService`](examples/command-line-only/internal/interfaces/email.go), [`CacheService`](examples/command-line-only/internal/interfaces/cache.go)) and a complete [`UserService`](examples/command-line-only/internal/service/user_service.go) with business logic, the results were clear.

The direct tooling approach won decisively, with the Makefile variant taking the crown:

**The Real Comparison**:

- **50 lines of Makefile** vs. **1000+ lines of Go code** for the MCP server
- **Zero infrastructure** vs. **Docker deployment and JSON-RPC**
- **Direct command execution** vs. **Complex protocol implementation**

This was honestly surprising. I expected the MCP server to be the clear winner, but simplicity trumped sophistication. The Makefile approach provides 90% of the benefits with 10% of the complexity.

Here's what made the difference:

- **Dead simple**: Everyone understands `make generate-mocks`
- **Self-documenting**: `make help` shows all available commands
- **Fast**: No server startup or protocol overhead
- **Reliable**: Uses proven Mockery with consistent configuration
- **Composable**: Easy to add new targets like `make test-coverage`

## What This Actually Means

This experiment shifted my thinking about AI-assisted development fundamentally. Instead of asking "How do I make AI write better code?", I started asking "How do I make AI better at using my existing tools?"

### AI as Orchestrator, Not Generator

The pattern that emerged is powerful: AI excels at understanding context, making decisions, and coordinating workflows. It's much less reliable at generating consistent, framework-compliant code from scratch.

### Simplicity Wins (Again)

The most sophisticated solution isn't always the best solution. The Makefile approach requires minimal setup but delivers maximum reliability. When it breaks, it breaks in predictable ways.

### Documentation Matters More Than Code

The [`CLAUDE.md`](examples/command-line-only/CLAUDE.md) files turned out to be as critical as the implementations themselves. AI needs clear, specific instructions. Vague guidelines produce vague results.

### The Bigger Pattern

This "AI as tool orchestrator" pattern likely extends far beyond mock generation:

- Database migration generation
- API client generation
- Infrastructure as Code templates
- Documentation generation
- Code refactoring workflows

The key insight: AI doesn't need to be good at everything. It just needs to understand what you want and know which tools to use to get there.

## Try It Yourself

All the code in this repository is functional and tested. No prototypes or proof-of-concepts here. The Makefile approach is particularly easy to experiment with:

```bash
cd examples/makefile
make help          # See all available commands
make generate-mocks # Generate all mocks
make test          # Run tests with generated mocks
```

Or try the command-line version:

```bash
cd examples/command-line-only
go mod tidy
# Follow the instructions in CLAUDE.md with your favorite AI assistant
```

## The Takeaway

When Mockery adds new features, my AI workflows automatically benefit. No prompt engineering required. No complex infrastructure to maintain. Just reliable, predictable mock generation that actually works.

The best mock, it turns out, is the one that actually works. And sometimes the best way to get there is to stop fighting your tools and start orchestrating them instead.

---

*Pick the approach that fits your team's complexity tolerance. The approaches work, the comparisons are honest, and the learnings are real.*
