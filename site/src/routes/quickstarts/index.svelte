<script context="module">
	export async function preload() {
        const qss = await this.fetch(`quickstarts.json`).then(r => r.json());
        console.log(qss);
        console.log('quickstarts: '+qss.length);
		return { qss };
	}
</script>

<script>
	const description = "Quickstarts ZITADEL";
	export let qss;
</script>

<svelte:head>
	<title>Quickstarts ZITADEL</title>

	<meta name="twitter:title" content="Quickstarts ZITADEL">
	<meta name="twitter:description" content={description}>
	<meta name="Description" content={description}>
</svelte:head>

<div class='quickstarts stretch'>
	<h1>Quickstarts</h1>
	{#each qss as qs}

		<article class='quickstart'>
            <div>
                <p class="sub">{qs.metadata.subtitle}</p>
                <a class="anchor" sapper:prefetch href={qs.fragment} title='{qs.metadata.title}'>
                    <h2>
                        <span id={qs.fragment} class="offset-anchor"></span>
                        {qs.metadata.title}
                    </h2>
                </a>
                <p>{@html qs.description}</p>
                <p class="info">{qs.metadata.date} • {qs.metadata.readingtime}</p>
            </div>
            <img src={qs.metadata.img} alt="article img" />
		</article>
	{/each}
	<p class="disclaimer">See also our Github page <a href="https://github.com/caos/zitadel" rel="external">ZITADEL </a> for questions regarding the sourcecode.</p>
</div>

<style>
	.quickstarts {
		grid-template-columns: 1fr 1fr;
		grid-gap: 1em;
		min-height: calc(100vh - var(--nav-h));
		padding: var(--top-offset) var(--side-nav) 6rem var(--side-nav);
		max-width: var(--main-width);
        margin: 0 auto 0 auto;
	}
	h2 {
		display: inline-block;
		color: white;
		max-width: 18em;
		font-size: var(--h3);
        font-weight: 400;
        margin: .5rem 0;
    }
    
    .quickstart {
        margin: 3.2rem 0 1rem 0;
        display: flex;
        flex-wrap: wrap;
        align-items: center;
    }

    .quickstart img {
        max-width: 150px;
        object-fit: contain;
        padding: 20px;
    }

    .quickstart .sub {
        font-size: 12px;
        font-size: bold;
        text-transform: uppercase;
        color: var(--second);
        margin: 0;
    }

    .quickstart .info {
        font-size: 12px;
        font-size: bold;
        text-transform: uppercase;
        color: var(--dark-text);
        margin: 0;
    }

	.quickstart:first-child {
		margin: 0 0 2rem 0;
		padding: 0 0 4rem 0;
		border-bottom: var(--border-w) solid #6767785b; /* based on --second */
	}
	.quickstart:first-child h2 {
		font-size: 4rem;
		font-weight: 400;
		color: var(--second);
    }

	.quickstart p {
		font-size: var(--h5);
		max-width: 30em;
		color: var(--dark-text);
	}
	:global(.quickstarts .quickstart ul) {
		margin-left: 3.2rem;
	}
	.quickstarts :global(.anchor) {
		top: calc((var(--h3) - 24px) / 2);
    }
    
    .quickstart:last-child {
        margin-bottom: 100px;
    }
	@media (max-width: 768px) {
		.quickstarts :global(.anchor) {
			transform: scale(0.6);
			opacity: 1;
			left: -1.0em;
		}
    }
    
    .disclaimer {
        margin-top: 100px;
    }
</style>