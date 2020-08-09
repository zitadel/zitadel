<script context="module">
    export async function preload(page, session) {
        const {params} = page;
        let workflow = '';
        try {
            workflow = `name: Docs
on:
  push:
    branches:
      - master

env:
  GITHUB_TOKEN: \${{ secrets.GITHUB_TOKEN }}

jobs:
  builddocs:
    name: Build Doc Frontend
    runs-on: ubuntu-latest
    steps:
    - name: Checkout Repo
      uses: actions/checkout@master
    - name: Build Docs
      uses: caos/site@master
      with:
        args: --basepath \${{ github.event.repository.name }}
    - name: Archive Production Artifact
      uses: actions/upload-artifact@master
      with:
        name: export
        path: __sapper__/export/\${{ github.event.repository.name }}
  deploydocs:
    name: Deploy
    needs: builddocs
    runs-on: ubuntu-latest
    steps:
      - name: Checkout Repo
        uses: actions/checkout@master
      - name: Download Artifact
        uses: actions/download-artifact@master
        with:
          name: export
          path: __sapper__/export/\${{ github.event.repository.name }}
      - name: Deploy
        uses: JamesIves/github-pages-deploy-action@releases/v3
        with:
          ACCESS_TOKEN: \${{ secrets.ACCESS_TOKEN }}
          BRANCH: gh-pages
          FOLDER: __sapper__/export/\${{ github.event.repository.name }}
          CLEAN: true
`;//await this.fetch(`.github/workflows/main.yml`, 'utf-8').then(x => x.text());
        } catch (err) {
            console.error(err);
            workflow = '';
        }
        return {workflow};
    }
</script>

<script>
    import Split from "../components/Split.svelte";
    import Section from '../components/Section.svelte';
    export let workflow;
    import { _ } from 'svelte-i18n';
</script>

<style>
    .caos-back {
        position: absolute;
        top: 0;
        right: 0;
        height: 70vh;
    }

    .logo {
        width: 100px;
        margin-bottom: 1rem;
    }
</style>

<svelte:head>
  <title>
    Caos • Documentation
  </title>
</svelte:head>

<img class="caos-back" src="logos/caos-logo-outline-1px-pictureonly.svg" alt="caos logo">

<Section>
<!-- <CodeTable></CodeTable> -->
</Section>
<Section>
<img class="logo" src="icons/android-chrome-256x256.png" alt="logo">
    <h2>{$_('title')}</h2>
    <p>{$_('description')}<a href="/get_started" >{$_('description_link')}</a></p>
<Split>
    <div class="description" slot="what">
        <p>
       {$_('instruction_header_0')}
        </p>

        <p>{$_('instruction_header_1')}</p>
    </div>

    <div
        style="grid-area: start; display: flex; flex-direction: column; min-width: 0"
        slot="how">
        <pre
        class="language-bash"
        style="margin: 0 0 1em 0; min-width: 0; min-height: 0">
        {workflow}
        </pre>
    </div>
</Split>
</Section>
