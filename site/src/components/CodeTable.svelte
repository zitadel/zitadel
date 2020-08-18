<script>
    let tabs = [
        {lang: 'yaml', content:"main: halllo"},
        {lang: 'js', content:"console.log();"},
        {lang: 'html', content:"<p>hello</p>"},
    ]
</script>

<style>
.pc-tab > input,
.pc-tab section > div {
  display: none;
}

#tab1:checked ~ section .tab1,
#tab2:checked ~ section .tab2,
#tab3:checked ~ section .tab3 {
  display: block;
}

#tab1:checked ~ nav .tab1,
#tab2:checked ~ nav .tab2,
#tab3:checked ~ nav .tab3 {
  color: red;
}

/* Visual Styles */
*, *:after, *:before {
  -webkit-box-sizing: border-box;
  -moz-box-sizing: border-box;
  box-sizing: border-box;
}

.pc-tab {
  width: 100%;
  max-width: 700px;
  margin: 0 auto;
}
.pc-tab ul {
  list-style: none;
  margin: 0;
  padding: 0;
}
.pc-tab ul li label {
  float: left;
  padding: 10px 25px;
  border: 1px solid #8795a1;
  border-bottom: 0;
  background: #212224;
}
.pc-tab ul li:first-child label {
  border-top-left-radius: 8px;
}

.pc-tab ul li:last-child label {
  border-top-right-radius: 8px;
}
.pc-tab ul li label:hover {
  background: #ffffff20;
}
.pc-tab ul li label:active {
  background: #ffffff30;
}
.pc-tab ul li:not(:last-child) label {
  border-right-width: 0;
}
.pc-tab section {
  clear: both;
}
.pc-tab section div {
  padding: 20px;
  width: 100%;
  border: 1px solid #8795a1;
  background: #1e1e1e;
  line-height: 1.5em;
  letter-spacing: 0.3px;
}

#tab1:checked ~ nav .tab1 label,
#tab2:checked ~ nav .tab2 label,
#tab3:checked ~ nav .tab3 label {
  background: #1e1e1e;
  position: relative;
  color: white;
}
#tab1:checked ~ nav .tab1 label:after,
#tab2:checked ~ nav .tab2 label:after,
#tab3:checked ~ nav .tab3 label:after {
  content: '';
  display: block;
  position: absolute;
  height: 2px;
  width: 100%;
  background: #1e1e1e;
  left: 0;
  bottom: -1px;
}

</style>

<div class="pc-tab">
<input checked id="tab1" type="radio" name="pct" />
<input id="tab2" type="radio" name="pct" />
<input id="tab3" type="radio" name="pct" />
  <nav>
    <ul>
        {#each tabs as { lang }, i}
        <li class="tab{i+1}">
            <label for="tab{i+1}">{lang}</label>
        </li>
        {/each}
    </ul>
  </nav>
  <section>
   {#each tabs as { lang, content }, i}
        <div class="tab{i+1}">
            {content}
        </div>
    {/each}
  </section>
</div>
