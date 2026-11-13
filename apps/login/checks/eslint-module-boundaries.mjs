// Keeps Server Actions and server-only code in separate, recognizable places.
//
// - src/lib/actions/: only Server Actions. Every file starts with "use server", and every export is a
//   public, unauthenticated endpoint, so its arguments must be treated as attacker-controlled.
// - src/lib/server/: server-only code. Every file imports "server-only" and none is a Server Action.
// - "use server" is not allowed anywhere else, neither as a file directive nor inline in a function.
// - Client components must not import from src/lib/server/.
//
// The build additionally checks the published actions against checks/server-actions.allowlist.json.

const ACTIONS_DIR = "/src/lib/actions/";
const SERVER_DIR = "/src/lib/server/";
const SERVER_IMPORT = /^@\/lib\/server(\/|$)/;

function directivesOf(program) {
  return program.body
    .filter((statement) => statement.type === "ExpressionStatement" && typeof statement.directive === "string")
    .map((statement) => statement.directive);
}

const moduleBoundaries = {
  meta: {
    type: "problem",
    docs: { description: "Keep Server Actions in src/lib/actions and server-only code in src/lib/server" },
    messages: {
      useServerOutsideActions:
        '"use server" is only allowed in src/lib/actions/. Every export of such a module is a public endpoint; put the action there and the implementation in src/lib/server/.',
      missingUseServer: 'Files in src/lib/actions/ must start with "use server".',
      missingServerOnly: 'Files in src/lib/server/ must import "server-only".',
      serverImportInClient:
        "Client components must not import from src/lib/server/. Call a Server Action from src/lib/actions/ instead.",
    },
    schema: [],
  },
  create(context) {
    const filename = context.filename.replaceAll("\\", "/");
    const isTest = /\.test\.tsx?$/.test(filename);
    const inActions = filename.includes(ACTIONS_DIR);
    const inServer = filename.includes(SERVER_DIR);

    return {
      Program(program) {
        if (isTest) {
          return;
        }

        const directives = directivesOf(program);

        if (inActions && !directives.includes("use server")) {
          context.report({ node: program, messageId: "missingUseServer" });
        }

        if (
          inServer &&
          !program.body.some((statement) => statement.type === "ImportDeclaration" && statement.source.value === "server-only")
        ) {
          context.report({ node: program, messageId: "missingServerOnly" });
        }

        if (directives.includes("use client")) {
          for (const statement of program.body) {
            if (statement.type === "ImportDeclaration" && SERVER_IMPORT.test(statement.source.value)) {
              context.report({ node: statement, messageId: "serverImportInClient" });
            }
          }
        }
      },
      // file-level and function-level directives are both ExpressionStatements with a `directive`
      "ExpressionStatement[directive='use server']"(node) {
        if (!inActions) {
          context.report({ node, messageId: "useServerOutsideActions" });
        }
      },
    };
  },
};

export default {
  rules: {
    "module-boundaries": moduleBoundaries,
  },
};
