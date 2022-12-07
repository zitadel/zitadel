import { Section } from "../../components/Section";

export const section = {
  render: Section,
  children: ["paragraph", "tag", "list", "code"],
  attributes: {
    columns: {
      type: Number,
      default: 1,
    },
  },
};
