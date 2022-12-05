import { nodes } from "@markdoc/markdoc";
import { CodeBlock } from "../../components/CodeBlock";

export const fence = {
  render: CodeBlock,
  attributes: nodes.fence.attributes,
};
