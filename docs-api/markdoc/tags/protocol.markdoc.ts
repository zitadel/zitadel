import type { Node, Config } from "@markdoc/markdoc";
import { Protocol } from "../../components/Protocol";

export default {
  render: Protocol,
  children: [""],
  attributes: {
    showDefault: {
      type: Boolean,
    },
    showOnProtocol: {
      type: String,
    },
  },
};
