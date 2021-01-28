<script>
	import { afterUpdate } from 'svelte';
	import Icon from './Icon.svelte';
    import CodeTable from './CodeTable.svelte';
    export let dir = '';
    export let language = 'en';
	export let sections = [];
	export let active_section = null;
	export let show_contents;
    export let prevent_sidebar_scroll = false;

	let ul;

	afterUpdate(() => {
		// bit of a hack — prevent sidebar scrolling if
		// TOC is open on mobile, or scroll came from within sidebar
		if (prevent_sidebar_scroll || show_contents && window.innerWidth < 832) return;

		const active = ul.querySelector('.active');

		if (active) {
			const { top, bottom } = active.getBoundingClientRect();

			const min = 200;
			const max = window.innerHeight - 200;

			if (top > max) {
				ul.parentNode.scrollBy({
					top: top - max,
					left: 0,
					behavior: 'smooth'
				});
			} else if (bottom < min) {
				ul.parentNode.scrollBy({
					top: bottom - min,
					left: 0,
					behavior: 'smooth'
				});
			}
		}
	});
</script>

<style>
	.reference-toc li {
		display: block;
		line-height: 1.2;
        margin: 0 0 4rem 0;
	}

	a {
		position: relative;
		transition: color 0.2s;
		border-bottom: none;
		padding: 0;
		color: var(--dark-text);
	}

	.section {
		display: block;
		padding: .8rem 0 .8rem 0;
		font-size: var(--h6);
		text-transform: uppercase;
		letter-spacing: 0.1em;
		font-weight: 600;
	}

	.subsection {
        display: block;
        font-family: var(--font);
		font-size: 14px;
        padding: 0.4em 0 0.4em 0;
        color: #a3acb9;
    }
    
    .section,
    .subsection {
        border-top-right-radius: 50vw;
        border-bottom-right-radius: 50vw;
        padding-left: 2rem;
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

	.section:hover,
	.subsection:hover,
	.active {
        color: var(--flash);
        font-weight: 500;
        color: #6c8eef;
        background-color: rgba(82,130,193,.1);
        padding-right: 1rem;
	}

	.subsection[data-level="4"] {
		padding-left: 3.2rem;
	}

	.icon-container {
        margin-left: .5rem;
        color: white;
    }

	@media (min-width: 832px) {
		a {
			color: var(--sidebar-text);
		}

		a:hover,
		.section:hover,
		.subsection:hover,
		.active {
			color: var(--prime);
		}
    }
</style>

<ul
	bind:this={ul}
	class="reference-toc"
	on:mouseenter="{() => prevent_sidebar_scroll = true}"
	on:mouseleave="{() => prevent_sidebar_scroll = false}"
>
	{#each sections as section}
		<li>
			<a class="section" class:active="{section.slug === active_section}" href="{language}/{dir}#{section.slug}">
				{@html section.metadata.title}

				{#if section.slug === active_section}
					<div class="icon-container">
						<Icon name="las la-arrow-right" />
					</div>
				{/if}
			</a>

			{#each section.subsections as subsection}
				<!-- see <script> below: on:click='scrollTo(event, subsection.slug)' -->
				<a
					class="subsection"
					class:active="{subsection.slug === active_section}"
					href="{language}/{dir}#{subsection.slug}"
					data-level="{subsection.level}"
				>
					{@html subsection.title}

					{#if subsection.slug === active_section}
						<div class="icon-container">
							<Icon name="las la-arrow-right" />
						</div>
					{/if}
				</a>
			{/each}
		</li>
	{/each}
</ul>
