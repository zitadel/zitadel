<style>
	.split {
		display: grid;
		grid-row-gap: 1em;
		grid-template-areas:
			"one"
			"two"
			"three"
			"what"
			"how";
	}

	.split :global(p) {
		font-size: var(--h5);
    }

	.how {
		/* needed to prevent the <pre> from
		   breaking the grid layout */
		min-width: 0;
	}

	.how :global(.cta) :global(a) {
		display: inline-block;
		text-align: right;
		background-color: var(--prime);
		padding: 0.5em 1.8em 0.5em 1em;
		border-radius: 16px;
		color: white;
		position: relative;
	}

	.how :global(.cta) :global(a)::after {
		right: 0.5em;
		top: 0.75em;
	}

	.what {
		margin: 2em 0 0 0;
	}

	.how :global(.cta) {
		margin: 0;
	}

	@media (min-width: 900px) {
		.split {
			grid-column-gap: 1em;
			grid-row-gap: 1em;
			grid-template-columns: repeat(2, 1fr);
			grid-template-areas:
				"one two"
				"three how"
				"what what";
		}
	}

	@media (min-width: 1200px) {
		.split {
			grid-column-gap: 1em;
			grid-row-gap: 5em;
			grid-template-columns: repeat(6, 1fr);
			grid-template-areas:
				"one one two two three three"
				"what what what how how how";
		}

		.what {
			margin: 0;
		}
	}
</style>

<div class="split">
    <div class="what" style="grid-area: what;">
        <slot name="what"></slot>
    </div>

    <div class="how" style="grid-area: how;">
        <slot name="how"></slot>
    </div>
</div>
