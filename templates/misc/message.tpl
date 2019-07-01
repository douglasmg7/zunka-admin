{{template "base" .}}

{{ define "embedded-css"}}
<style type="text/css">
    .content {
        max-width: var(--width-small);
    }
</style>
{{end}}

{{define "title"}}Message{{end}}

{{define "content"}}
<div class="content">
    <!-- Message title. -->
    <h2>{{.TitleMsg}}</h2>

    <!-- Foot messages. -->
    {{if .SuccessMsg}} <div class="success-msg"> {{.SuccessMsg}} </div> {{end}}
    {{if .WarnMsg}} <div class="warn-msg"> {{.WarnMsg}} </div> {{end}}
</div>
{{end}}