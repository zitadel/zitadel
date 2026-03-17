import { Callout } from 'fumadocs-ui/components/callout';

interface TermNoteProps {
    newTerm: string;
    oldTerms: string[];
}

export function TerminologyUpdate({ newTerm, oldTerms }: TermNoteProps) {
    // Formats the list: "A, B and C"
    const formattedOldTerms = oldTerms.map((term, i) => (
        <b key={term}>
            {term}{i < oldTerms.length - 2 ? ', ' : i === oldTerms.length - 2 ? ' or ' : ''}
        </b>
    ));

    return (
        <Callout type="info">
            <b>Terminology Update:</b> We have streamlined our naming conventions to improve clarity.
            The term <b>{newTerm}</b> now replaces what was previously referred to as {formattedOldTerms}.
            These terms all refer to the same underlying functionality.
        </Callout>
    );
}