import { Column } from "../../components/Column";

export default {
  render: Column,
  children: ["paragraph", "tag", "list", "code"],
  attributes: {
    position: {
      type: String,
    },
  },
};
