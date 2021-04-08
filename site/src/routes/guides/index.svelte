<script context="module">
	export async function preload() {
        const qss = await this.fetch(`guides.json`).then(r => r.json());
        return { qss };
	}
</script>

<script>
	const description = "Guides ZITADEL";
	export let qss;
</script>

<svelte:head>
	<title>ZITADEL Guides</title>

	<meta name="twitter:title" content="Guides ZITADEL">
	<meta name="twitter:description" content={description}>
	<meta name="Description" content={description}>
</svelte:head>

<div class='guides stretch'>
	<h1>Guides ZITADEL</h1>
	{#each qss as qs}
        {#if qs.metadata.visible == 'true'} 
        <article class='guide'>
            <div>
                <p class="sub">{qs.metadata.subtitle}</p>
                <h2>
                    <span id={qs.fragment} class="offset-anchor"></span>
                    <a class="anchor" sapper:prefetch href='guides#{qs.fragment}' title='{qs.title}'><i class="las la-link"></i></a>
                    {qs.metadata.title}
                </h2>
                <p>{@html qs.answer}</p>
                <a class="link" href="{qs.fragment}" sapper:prefetch>Read Guide <i class="las la-arrow-right"></i></a>
                <p class="info">{qs.metadata.date} • {qs.metadata.readingtime}</p>
            </div>
            <img src={qs.metadata.img} alt="article img" />
        </article>
        {/if}
	{/each}
<p class="disclaimer">See also our Github page <a href="https://github.com/caos/zitadel" rel="external">ZITADEL </a> for questions regarding the sourcecode.</p></div>

<style>
	.guides {
		grid-template-columns: 1fr 1fr;
		grid-gap: 1em;
		min-height: calc(100vh - var(--nav-h));
		padding: var(--top-offset) var(--side-nav) 6rem var(--side-nav);
		max-width: var(--main-width);
        margin: 0 auto 0 auto;
    }
    
	h2 {
        position: relative;
		display: inline-block;
		color: white;
		max-width: 18em;
		font-size: var(--h3);
        font-weight: 400;
        margin: .5rem 0;
    }
    
    .guide {
        margin: 3.2rem 0 1rem 0;
        display: flex;
        flex-wrap: wrap;
        align-items: center;
    }

    .guide .anchor {
        position: absolute;
        display: block;
        background-size: 30px 30px;
        width: 30px;
        height: 30px;
        left: -1.3em;
        opacity: 0;
        color: white;
        transition: opacity 0.2s;
        border: none !important;
    }

    h2:hover .anchor {
        opacity: 1;
    }

    .guide img {
        max-width: 150px;
        max-height: 150px;
        object-fit: contain;
        padding: 20px;
        border-radius: 4px;
    }

    .guide .sub {
        font-size: 12px;
        font-size: bold;
        text-transform: uppercase;
        color: var(--second);
        margin: 0;
    }

    .guide .info {
        font-size: 12px;
        font-size: bold;
        text-transform: uppercase;
        color: var(--dark-text);
        margin: 1rem 0 0 0;
    }

    .guide .link,
    .guide .link i{
        color: var(--prime);
        text-decoration: none;
    }

	.guide:first-child {
		margin: 0 0 2rem 0;
		padding: 0 0 4rem 0;
		border-bottom: var(--border-w) solid #6767785b; /* based on --second */
	}
	.guide:first-child h2 {
		font-size: 4rem;
		font-weight: 400;
		color: var(--second);
    }

	.guide p {
		font-size: var(--h5);
		max-width: 30em;
		color: var(--dark-text);
	}
	:global(.guides .guide ul) {
		margin-left: 3.2rem;
	}
	.guides :global(.anchor) {
		top: calc((var(--h3) - 24px) / 2);
    }
    
    .guide:last-child {
        margin-bottom: 100px;
    }

	@media (max-width: 768px) {
        .guide {
            flex-direction: column-reverse;
            margin-bottom: 50px;
        }

        .guide img {
            max-height: 180px;
            max-width: 180px;
            padding-bottom: 10px;
        }

		.guides :global(.anchor) {
			transform: scale(0.6);
			opacity: 1;
			left: -1.0em;
		}
    }
    
    .disclaimer {
        margin-top: 100px;
    }
</style>