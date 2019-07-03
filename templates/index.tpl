{{ define "embedded-css"}}
<style type="text/css">
    .img-title {
        /* margin-bottom: .1em; */
        /* color: #3878BB; */
    }
    h2, h3 {
        margin-top: 0;
    }
    .panel {
        margin-top: 1em;
    }
    img {
        border-radius: 4px;
    }
</style>
{{end}}

{{template "base" .}}

{{define "title"}}{{end}}

{{define "content"}}
<div class="content">
    <div>
        <h1></h1>
        <h3></h3>
    </div>
</div>
{{end}}
