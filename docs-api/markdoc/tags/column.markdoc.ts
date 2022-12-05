import { Column } from "../../components/Column";

export const column = {
  render: Column,
  children: ["paragraph", "tag", "list", "code"],
  attributes: {
    title: {
      type: String,
    },
    language: {
      type: String,
    },
  },
};
