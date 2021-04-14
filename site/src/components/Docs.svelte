<script>
  import { onMount } from "svelte";
  import GuideContents from "./GuideContents.svelte";
  import Icon from "./Icon.svelte";
  export let owner = "caos";
  export let path = "docs";
  export let project = "zitadel";	
  export let dir = "";
  export let edit_title = "edit this section";
  export let sections;
  import SearchSelector from './SearchSelector.svelte';
  import SearchTrigger from './SearchTrigger.svelte';
  let searchEnabled = false;
  let active_section;

  let container;
  let aside;
  let show_contents = false;

  function handleSearch(event) {
    searchEnabled = !event.detail.closed;
  }

  function handleKeydown(event) {
      const isCtrlKey = navigator.platform.indexOf('Mac') > -1 ? event.metaKey : event.ctrlKey;
      const isShiftKey = event.shiftKey;
    if ((event.keyCode == 114 && isShiftKey) || (isCtrlKey && event.keyCode == 70 && isShiftKey)) {
        event.preventDefault();
        searchEnabled = !searchEnabled;
    }
  }

  onMount(() => {
    // don't update `active_section` for headings above level 4, see _sections.js
    const anchors = container.querySelectorAll("[id]:not([data-scrollignore])");

    let positions;

    const onresize = () => {
      const { top } = container.getBoundingClientRect();
      positions = [].map.call(anchors, anchor => {
        return anchor.getBoundingClientRect().top - top;
      });
    };

    let last_id = window.location.hash.slice(1);

    const onscroll = () => {
      const top = -window.scrollY;

      let i = anchors.length;
      while (i--) {
        if (positions[i] + top < 40) {
          const anchor = anchors[i];
          const { id } = anchor;

          if (id !== last_id) {
            active_section = id;
            last_id = id;
          }

          return;
        }
      }
    };

    window.addEventListener("scroll", onscroll, true);
    window.addEventListener("resize", onresize, true);

    // wait for fonts to load...
    const timeouts = [setTimeout(onresize, 1000), setTimeout(onscroll, 5000)];

    onresize();
    onscroll();

    return () => {
      window.removeEventListener("scroll", onscroll, true);
      window.removeEventListener("resize", onresize, true);

      timeouts.forEach(timeout => clearTimeout(timeout));
    };
  });
</script>

<style>
  .overlay {
      position: fixed;
      top: var(--nav-h);
      right: 0;
      left: 0;
      bottom: 0;
      background: #00000050;
      backdrop-filter: blur(10px);
      visibility: hidden;
  }

  .overlay.visible {
      visibility: visible;
  }

  aside {
    position: fixed;
    background-color: var(--side-nav-back);
    left: 0.8rem;
    bottom: 0.8rem;
    width: 3.4rem; /* size to match button */
    height: 3.4rem; /* size to match button */
    border-radius: .5rem;
    overflow: hidden;
    border: 1px solid #8795a1;
    box-shadow: 1px 1px 6px rgba(0, 0, 0, 0.1);
    transition: width 0.2s, height 0.2s;
  }

  aside button {
    position: absolute;
    bottom: 0;
    left: 0;
    width: 3.4rem;
    height: 3.4rem;
  }

  aside.open {
    width: calc(100vw - 1.5rem);
    height: calc(100vh - var(--nav-h) - 15rem);
  }

  aside.open::before {
    content: "";
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 2em;
    pointer-events: none;
    z-index: 2;
  }

  .sidebar {
    position: absolute;
    font-family: var(--font);
    overflow-y: auto;
    width: 100%;
    height: 100%;
    padding: 4em 1.6rem 2em 0;
    bottom: 2em;
  }

  .sidebar :global(.language-switcher) {
        position: relative;
  }

  aside .sidebar :global(.search-trigger) {
      visibility: hidden;
  }

  aside.open .sidebar :global(.search-trigger) {
      visibility: visible;
  }

  .sidebar::-webkit-scrollbar-track {
    -webkit-box-shadow: inset 0 0 6px rgba(0, 0, 0, .3);
    box-shadow: inset 0 0 6px rgba(0, 0, 0, .3);
    background-color: #00000010;
    border-radius: 8px;
  }

  .sidebar::-webkit-scrollbar {
    width: 6px;
    height: 6px;
    background-color: #00000010;
  }

  .sidebar::-webkit-scrollbar-thumb {
    background-color: #6c8eef30;
    border-radius: 8px;
    cursor: pointer;
  }

  .content {
    width: 100%;
    margin: 0;
    padding: calc(var(--top-offset) + 100px) var(--side-nav);
    tab-size: 2;
    -moz-tab-size: 2;
  }

  @media (min-width: 832px) {
    /* can't use vars in @media :( */
    aside {
      display: block;
      width: var(--sidebar-w);
      height: calc(100vh - var(--searchbar-space));
      top: var(--searchbar-space); /* space for searchbar */
      left: 0;
      overflow: hidden;
      box-shadow: none;
      border: none;
      overflow: hidden;
      background-color: var(--side-nav-back);
      color: white;
    }

    aside.open::before {
      display: none;
    }

    aside button {
      display: none;
    }

    aside .sidebar :global(.search-trigger) {
      visibility: visible;
    }

    .sidebar {
      padding: 1em .5rem 6.4rem 0;
      font-family: var(--font);
      overflow-y: auto;
      height: 100%;
      bottom: auto;
      width: 100%;
    }

    .sidebar :global(.language-switcher) {
        display: flex;
        position: fixed;
    }

    .sidebar :global(.search-trigger) {
      visibility: visible;
    }

    .content {
      padding-left: calc(var(--sidebar-w) + var(--side-nav));
    }

    .content :global(.side-by-side) {
      display: grid;
      grid-template-columns: calc(50% - 0.5em) calc(50% - 0.5em);
      grid-gap: 1em;
    }

    .content :global(.side-by-side) :global(.code) {
      padding: 1em 0;
    }
  }

  .content h2 {
    margin-top: 8rem;
    padding: 2rem 1.6rem 4rem 0.2rem;
    border-top: var(--border-w) solid #6767785b; /* based on --second */
    color: var(--heading);
    line-height: 1;
    font-size: var(--h3);
    letter-spacing: 0.05em;
    text-transform: uppercase;
  }

  .content section:first-of-type > h2 {
    margin-top: -8rem;
    border-top: none;
  }

  .content :global(h4) {
    margin: 2em 0 1em 0;
  }

  .content :global(.offset-anchor) {
    position: relative;
    display: block;
    top: calc(-1 * (var(--nav-h) + var(--top-offset) - 1rem));
    width: 0;
    height: 0;
  }

  .content :global(.anchor) {
    position: absolute;
    display: block;
    /* TODO replace link icon */
    /* background: url(../icons/link.svg) 0 50% no-repeat; */
    background-size: 30px 30px;
    width: 30px;
    height: 30px;
    left: -1.3em;
    opacity: 0;
    color: white;
    transition: opacity 0.2s;
    border: none !important; /* TODO get rid of linkify */
  }

  .content :global(.anchor) :global(i) {
    color: #8795a1;
    /* font-size: 24px; */
  }

  @media (min-width: 768px) {
    .content :global(h2):hover :global(.anchor),
    .content :global(h3):hover :global(.anchor),
    .content :global(h4):hover :global(.anchor),
    .content :global(h5):hover :global(.anchor),
    .content :global(h6):hover :global(.anchor) {
      opacity: 1;
    }

    .content :global(h5) :global(.anchor),
    .content :global(h6) :global(.anchor) {
      top: 0.2em;
    }
  }

  .content :global(h3),
  .content :global(h3 > code) {
    margin: 2.0rem 0 0 0;
    padding: 2rem 1.6rem 2.0rem 0.2rem;
    color: var(--heading);
    border-top: var(--border-w) solid #6767781f; /* based on --second */
    background: transparent;
    line-height: 1;
  }

  .content :global(h3):first-of-type {
    border: none;
    margin: 0;
  }

  /* avoid doubled border-top */
  .content :global(h3 > code) {
    border-radius: 0 0 0 0;
    border: none;
    font-size: inherit;
        font-weight: 700;
  }

  .content :global(h4),
  .content :global(h4 > code) {
    font-family: inherit;
    font-weight: 500;
    font-size: 2.4rem;
    color: var(--heading);
    margin: 2.0rem 0 1.6rem 0;
    padding-left: 0;
    background: transparent;
    line-height: 1;
    padding: 0;
    top: 0;
  }

  .content :global(h4 > em) {
    opacity: 0.7;
  }

  .content :global(h5) {
    font-size: 2.4rem;
    margin: 2em 0 0.5em 0;
  }

  .content :global(code) {
    padding: 0.3rem 0.8rem 0.3rem;
    margin: 0 0.2rem;
    top: -0.1rem;
    background: #2a2f45;
  }

  .content :global(p) :global(code) {
    border: 1px solid #ffffff20;
  }

  .content :global(pre) :global(code) {
    padding: 0;
    margin: 0;
    top: 0;
    background: transparent;
  }

  .content :global(pre) {
    margin: 0 0 2em 0;
    width: 100%;
    max-width: 100%;
  }

  .content :global(.icon) {
    width: 2rem;
    height: 2rem;
    stroke: currentColor;
    stroke-width: 2;
    stroke-linecap: round;
    stroke-linejoin: round;
    fill: none;
  }

  .content :global(table) {
    margin: 0 0 2em 0;
  }

  section > :global(.code-block) > :global(pre) {
    display: inline-block;
    /* background: var(--back-api); */
    color: white;
    padding: 0.3rem 0.8rem;
    margin: 0;
    max-width: 100%;
  }

  section > :global(.code-block) > :global(pre.language-markup) {
    padding: 0.3rem 0.8rem 0.2rem;
    background: var(--back-api);
  }

  section > :global(p) {
    max-width: var(--linemax);
  }

  section :global(p) {
    margin: 1em 0;
    text-align: justify;
  }

  small {
    font-size: var(--h5);
    float: right;
    pointer-events: all;
    color: var(--prime);
    cursor: pointer;
  }

  /* no linkify on these */
  small a {
    all: unset;
  }

  small a:before {
    all: unset;
  }

  section :global(blockquote) {
    color: #85d996;
    border: 2px solid var(--grey-text);
    background: #2a2f45;
  }

  section :global(blockquote) :global(code) {
    /* background: hsl(204, 100%, 95%) !important; */
    color: var(--prime);
  }
</style>

<svelte:window on:keydown={handleKeydown}/>

<div bind:this={container} class="content listify">
  {#each sections as section}
    <section data-id={section.slug}>
      <h2>
        <span class="offset-anchor" id={section.slug} />
        <!-- svelte-ignore a11y-missing-content -->
        <a href="{dir}#{section.slug}" class="anchor" aria-hidden />

        {@html section.metadata.title}
        <small>
          <a
            href="https://github.com/{owner}/{project}/edit/main/site/{path}/{dir}/{section.file}"
            title={edit_title}>
            <Icon name="las la-external-link-alt" size="24px" />
          </a>
        </small>
      </h2>

      {@html section.html}
    </section>
  {/each}
</div>

<div class="overlay {show_contents ? 'visible' : ''}"></div>

<aside bind:this={aside} class="sidebar-container" class:open={show_contents}>
  <div class="sidebar" on:click={() => (show_contents = false)}>
    <SearchTrigger on:click={handleSearch}/>

    <!-- scroll container -->
    <GuideContents {dir} {sections} {active_section} {show_contents} />
  </div>

  <button on:click={() => (show_contents = !show_contents)}>
    <Icon name={show_contents ? 'las la-times' : 'las la-bars'} />
  </button>

</aside>

{#if searchEnabled == true}
    <SearchSelector on:close={handleSearch} {sections} slug={dir}></SearchSelector>
{/if}
