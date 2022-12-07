import { Column } from "../../components/Column";

export const column = {
  render: Column,
  children: ["paragraph", "tag", "list", "code"],
  attributes: {
    position: {
      type: String,
    },
  },
};
